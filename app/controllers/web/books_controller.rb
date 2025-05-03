class Web::BooksController < Web::ApplicationController
  before_action :authenticate_user!, only: [ :create_request, :download ]

  def show
    book_request = current_user.book_request

    render inertia: true, props: {
      bookRequested: !!book_request
    }
  end

  def create_request
    request = BookRequest.find_or_initialize_by user: current_user
    if request.new_record?
      request.save!
      event = BookRequestedEvent.new
      publish_event(event, current_user)
    end

    f(:success)
    redirect_to view_context.book_path
  end

  def download
    unless current_user.book_request
      redirect_to view_context.book_path
    end

    current_user.book_request.mark_as_downloaded!

    filepath = Rails.root.join("public/book.pdf")
    send_file filepath,
      filename: "profession-developer-hexlet-book.pdf",
      type: "application/pdf",
      disposition: "inline"
  end
end
