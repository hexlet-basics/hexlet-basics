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
    Rails.application.eager_load!
    events = ApplicationEvent.descendants.map { |klass| klass.name.gsub("::", "") }.sort

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
end
