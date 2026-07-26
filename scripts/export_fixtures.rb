# Export the legacy Rails test fixtures to Go testfixtures YAML in the repo-root
# `fixtures/` dir. Lives in the Go repo (NOT in legacy/); it only borrows the
# Rails runtime to resolve fixtures, run via:
#
#   cd legacy && RAILS_ENV=test bin/rails runner ../scripts/export_fixtures.rb
#
# (or just `make fixtures-import` from the repo root).
#
# Why the Rails runtime: resolving these fixtures (ERB, $LABEL, YAML merge keys,
# crc32 label ids, association -> *_id, namespaced tables, HABTM joins) IS the
# Rails fixture engine. We reuse it instead of reimplementing it, then dump the
# rows Rails would have inserted.
#
# Env: FORCE=1 overwrites existing files; ONLY=t1,t2 scopes to those tables.
# Non-destructive by default (skips files that already exist), so hand-curated
# fixtures are never clobbered.

require "fileutils"
require "yaml"

module FixtureExport
  module_function

  def run
    unless Rails.env.test?
      abort "run with RAILS_ENV=test (needs the test DB + rollback isolation)"
    end

    fixtures_dir = Rails.root.join("test/fixtures")
    out_dir      = File.expand_path("../fixtures", Rails.root) # repo-root/fixtures
    force        = ENV["FORCE"] == "1"
    only         = ENV["ONLY"]&.split(",")&.map(&:strip)
    internal     = %w[schema_migrations ar_internal_metadata]

    FileUtils.mkdir_p(out_dir)

    # Every *.yml under test/fixtures is a fixture set; its name is the path
    # relative to the dir, without extension (e.g. "language/categories").
    # `files/` is Rails' file_fixture_path (attachment blobs), not table
    # fixtures — exclude it, as fixture discovery does.
    set_names = Dir.glob(fixtures_dir.join("**/*.yml")).map do |path|
      Pathname.new(path).relative_path_from(fixtures_dir).to_s.delete_suffix(".yml")
    end.reject { |name| name.start_with?("files/") }.sort

    conn = ActiveRecord::Base.connection
    written = []
    skipped = []
    empty   = []

    # Resolve every fixture through Rails inside a transaction we roll back, so
    # the shared test DB is left untouched and Go test isolation is preserved.
    ActiveRecord::Base.transaction do
      ActiveRecord::FixtureSet.reset_cache
      ActiveRecord::FixtureSet.create_fixtures(fixtures_dir, set_names)

      tables = conn.tables.sort - internal
      tables.select! { |t| only.include?(t) } if only

      tables.each do |table|
        out_path = File.join(out_dir, "#{table}.yml")

        if File.exist?(out_path) && !force
          skipped << table
          next
        end

        rows = dump_rows(conn, table)
        if rows.empty?
          empty << table
          next
        end

        File.write(out_path, render_fixture(table, rows))
        written << table
      end

      raise ActiveRecord::Rollback
    end

    report(out_dir, written, skipped, empty)
  end

  # Read a table into ordered {column => value} hashes, casting each value to a
  # clean Ruby scalar so it serialises to readable YAML.
  def dump_rows(conn, table)
    quoted = conn.quote_table_name(table)
    pk     = conn.primary_key(table)
    order  = pk ? conn.quote_column_name(pk) : "1"
    result = conn.select_all("SELECT * FROM #{quoted} ORDER BY #{order}")

    result.map do |row|
      row.each_with_object({}) do |(col, raw), acc|
        acc[col] = cast_value(result.column_types[col], raw)
      end
    end
  end

  def cast_value(type, raw)
    return nil if raw.nil?

    value = type ? type.deserialize(raw) : raw
    case value
    when Time, DateTime then value.strftime("%Y-%m-%d %H:%M:%S")
    when Date           then value.strftime("%Y-%m-%d")
    when BigDecimal     then value.to_s("F")
    else value
    end
  end

  def render_fixture(table, rows)
    header = <<~HDR
      # Fixtures for `#{table}`, exported from the legacy Rails fixtures by
      # `make fixtures-import`. Loaded by testfixtures inside each test's rollback
      # transaction. Regenerate; do not hand-tune unless you stop regenerating.
    HDR
    # YAML.dump emits a sequence of mappings — exactly the testfixtures list form.
    header + YAML.dump(rows).sub(/\A---\n/, "")
  end

  def report(out_dir, written, skipped, empty)
    puts "Wrote #{written.size} fixture file(s) to #{out_dir}"
    puts "  written: #{written.join(', ')}" unless written.empty?
    puts "  skipped (already exist, use FORCE=1): #{skipped.join(', ')}" unless skipped.empty?
    puts "  empty (no rows, nothing to write): #{empty.join(', ')}" unless empty.empty?
    puts "\nGenerated ids are Rails' crc32 label ids (large): assert Go tests by " \
         "business fields (slug/email), not by raw id. Don't mix crc32-id " \
         "generated fixtures with small-id hand-curated ones in the same FK graph."
  end
end

FixtureExport.run
