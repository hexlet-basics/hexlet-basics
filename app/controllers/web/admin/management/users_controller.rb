# frozen_string_literal: true

class Web::Admin::Management::UsersController < Web::Admin::Management::ApplicationController
  # include ActionController::Live

  def index
    q = ransack_params("sf" => "id", "so" => "desc")
    search = User.ransack(q)
    pagy, records = pagy(search.result)

    render inertia: true, props: {
      users: UserResource.new(records),
      grid: GridResource.new(grid_params(pagy))
    }

    # respond_to do |format|
    #   format.html do
    #     @users = users.page(params[:page])
    #   end
    #
    #   format.csv do
    #     response.headers['Last-Modified'] = Time.now.httpdate.to_s
    #
    #     send_stream(filename: "users-#{params[:q][:language_members_created_at_gteq]}.csv") do |stream|
    #       stream.write "id, email, stack, finished_lessons\n"
    #
    #       users.find_each do |user|
    #         language_members = user.language_members
    #
    #         if language_members.any?
    #           language_members.each do |language_member|
    #             stream.write "#{user.id}, #{user.email}, #{language_member.language.name}, #{language_member.finished_lessons_count}\n"
    #           end
    #         else
    #           stream.write "#{user.id}, #{user.email}\n"
    #         end
    #       end
    #     end
    #   end
    # end
  end

  def edit
    user = Admin::UserForm.find params[:id]
    progressByLanguage = user.lesson_members.group(:language).count
    progress = progressByLanguage.map do |language, count|
      { language: language.slug, count: }
    end

    render inertia: true, props: {
      user: UserResource.new(user),
      progress:
    }
  end

  def update
    user = Admin::UserForm.find params[:id]

    if user.update(params[:user])
      f(:success)
    else
      f(:error)
    end

    redirect_to_inertia edit_admin_management_user_url(user), user
  end
end
