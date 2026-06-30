# typed: strict
# frozen_string_literal: true

class BlogPostRichTextContent
  extend T::Sig

  class << self
    extend T::Sig

    sig { params(html: String).returns(String) }
    def from_editor_html(html)
      fragment = Nokogiri::HTML5.fragment(html)

      fragment.css("img[data-blob-sgid]").each do |image_node|
        replace_editor_image_node!(image_node)
      end

      fragment.to_html
    end

    sig { params(html: String).returns(String) }
    def to_editor_html(html)
      fragment = Nokogiri::HTML5.fragment(html)

      fragment.css(ActionText::Attachment.tag_name).each do |attachment_node|
        replace_attachment_node!(attachment_node)
      end

      fragment.to_html
    end

    sig { params(html: String).returns(String) }
    def to_display_html(html)
      fragment = Nokogiri::HTML5.fragment(html)

      fragment.css(ActionText::Attachment.tag_name).each do |attachment_node|
        replace_attachment_node!(attachment_node)
      end

      fragment.to_html
    end

    private

    sig { params(image_node: T.untyped).void }
    def replace_editor_image_node!(image_node)
      sgid = image_node["data-blob-sgid"].to_s
      return if sgid.blank?

      blob = ActionText::Attachable.from_attachable_sgid(sgid)
      return unless blob.is_a?(ActiveStorage::Blob)

      attachment_html = ActionText::Attachment.from_attachable(blob).to_html
      image_node.replace(attachment_html)
    rescue ActiveRecord::RecordNotFound
      nil
    end

    sig { params(attachment_node: T.untyped).void }
    def replace_attachment_node!(attachment_node)
      attachable = ActionText::Attachable.from_node(attachment_node)
      return unless attachable.is_a?(ActiveStorage::Blob)

      image_node = Nokogiri::XML::Node.new("img", attachment_node.document)
      image_node["src"] = blob_path(attachable)
      image_node["alt"] = attachable.filename.to_s
      image_node["data-blob-sgid"] = attachable.attachable_sgid

      attachment_node.replace(image_node)
    rescue ActiveRecord::RecordNotFound
      nil
    end

    sig { params(blob: ActiveStorage::Blob).returns(String) }
    def blob_path(blob)
      Rails.application.routes.url_helpers.rails_service_blob_proxy_path(
        blob.signed_id,
        blob.filename,
      )
    end
  end
end
