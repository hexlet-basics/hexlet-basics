class ConvertBlogPostBodyToRichBody < ActiveRecord::Migration[8.1]
  COMMONMARKER_OPTIONS = {
    extension: { table: true, strikethrough: true, autolink: true, tasklist: true }
  }.freeze

  # One-time conversion of the legacy markdown `body` into Action Text `rich_body`.
  # The `::courses` directive is dropped — recommended courses now render as a
  # separate section on the blog post page.
  def up
    BlogPost.reset_column_information

    BlogPost.find_each do |post|
      markdown = post.read_attribute(:body).to_s
      next if markdown.blank?

      cleaned = markdown.each_line.reject { |line| line.strip == "::courses" }.join
      html = Commonmarker.to_html(cleaned, options: COMMONMARKER_OPTIONS)

      post.rich_body = html
      post.save(validate: false)
    end
  end

  def down
    raise ActiveRecord::IrreversibleMigration
  end
end
