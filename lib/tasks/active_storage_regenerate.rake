namespace :app do
  desc "Regenerate ActiveStorage variants"
  task regenerate_active_storage_variants: :environment do
    puts "=== Regenerating BlogPost cover variants ==="
    BlogPost.find_each do |post|
      next unless post.cover.attached?

      puts "Regenerating BlogPost ##{post.id} (#{post.slug})"
      begin
        post.cover.variant(:main).processed
        post.cover.variant(:list).processed
        post.cover.variant(:thumb).processed
      rescue => e
        warn "Failed to regenerate BlogPost ##{post.id}: #{e.message}"
      end
    end

    puts "\n=== Regenerating Language cover variants ==="
    Language.find_each do |lang|
      next unless lang.cover.attached?

      puts "Regenerating Language ##{lang.id} (#{lang.slug})"
      begin
        lang.cover.variant(:list).processed
        lang.cover.variant(:thumb).processed
      rescue => e
        warn "Failed to regenerate Language ##{lang.id}: #{e.message}"
      end
    end

    puts "\n=== Done regenerating all variants ==="
  end
end
