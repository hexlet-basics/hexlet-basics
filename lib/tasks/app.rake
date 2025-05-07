namespace :app do
  desc "Export event names to TypeScript file"
  task export_events_to_ts: :environment do
    registry = ApplicationContainer["event_registry"]
    events = registry.all.sort

    vite_root = ViteRuby.config.root # обычно app/frontend
    output_path = vite_root.join(ViteRuby.config.source_code_dir, "event_names.ts")

    File.open(output_path, "w") do |file|
      file.puts "// This file is auto-generated. Do not edit manually."
      file.puts "// Run `bin/rails app:export_events_to_ts` to update.\n\n"

      file.puts "export const eventNames = ["
      events.each do |event|
        file.puts "  \"#{event}\","
      end
      file.puts "] as const;"
    end

    puts "✅ Exported #{events.size} event names to #{output_path}"
  end
end
