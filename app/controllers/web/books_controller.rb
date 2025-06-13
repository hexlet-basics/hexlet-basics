class Web::BooksController < Web::ApplicationController
  before_action :authenticate_user!, only: [ :create_request, :download ]

  def show
    book_request = current_user.book_request

    render inertia: true, props: {
      lead: LeadCrudResource.new(LeadForm.new),
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
    book_request = current_user.book_request
    unless book_request
      redirect_to view_context.book_path
      return
    end

    book_request.state = "downloaded"
    book_request.save!

    filepath = Rails.root.join("public/book.pdf")
    send_file filepath,
      filename: "profession-developer-hexlet-book.pdf",
      type: "application/pdf",
      disposition: "inline"
  end
end
