require "fileutils"
require "json"

namespace :app do
  desc "Export event data types to TypeScript file"
  task export_events_to_ts: :environment do
    Rails.application.eager_load!
    output_path = Rails.root.join("app/javascript/types/events.ts")
    EventsTypesExporter.new(output_path: output_path).call
  end

  desc "Export event names to TypeScript file"
  task export_event_names_to_ts: :environment do
    event_registry = DepsLocator.current.event_registry
    events = event_registry.all.sort

    vite_root = ViteRuby.config.root # обычно app/frontend
    output_path = vite_root.join(ViteRuby.config.source_code_dir, "generated", "event_names.ts")

    File.open(output_path, "w") do |file|
      file.puts "// This file is auto-generated. Do not edit manually."
      file.puts "// Run `bin/rails app:export_events_to_ts` to update.\n\n"

      file.puts "export const eventNames = ["
      events.each do |event|
        file.puts "  '#{event}',"
      end
      file.puts "] as const;"
    end

    puts "✅ Exported #{events.size} event names to #{output_path}"
  end

  desc "Export Rails enums to TypeScript file"
  task export_enums_to_ts: :environment do
    Rails.application.eager_load!

    vite_root   = ViteRuby.config.root
    source_dir  = ViteRuby.config.source_code_dir
    out_dir     = vite_root.join(source_dir, "generated")
    output_path = out_dir.join("enums.ts")
    FileUtils.mkdir_p(out_dir)

    items = []

    ActiveRecord::Base.descendants.each do |model|
      next if model.abstract_class?
      begin
        next unless model.table_exists?
      rescue StandardError
        next
      end

      next unless model.respond_to?(:defined_enums)

      model.defined_enums.each do |enum_name, value_map|
        keys = value_map.keys
        next if keys.empty?

        # Admin::BlogPostForm.state → adminBlogPostFormState
        base = model.name.gsub("::", "").camelize(:lower)
        key  = "#{base}#{enum_name.camelize}"

        items << [ "#{model.name}.#{enum_name}", key, keys ]
      end
    end

    items.sort_by!(&:first)

    File.open(output_path, "w") do |file|
      file.puts "// This file is auto-generated. Do not edit manually."
      file.puts "// Run `bin/rails app:export_enums_to_ts` to update.\n\n"
      file.puts "/* eslint-disable */"
      file.puts "export const enums = {"

      items.each do |(_fq, key, keys)|
        file.puts "  #{key}: #{JSON.generate(keys)},"
      end

      file.puts "} as const;"
    end

    puts "✅ Exported #{items.size} enums to #{output_path}"
  end
end
