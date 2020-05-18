# frozen_string_literal: true

Rails.application.routes.draw do
  scope module: :web do
    root 'home#index'

    resource :session, only: %i[new create destroy]
    resources :registrations, only: %i[new create], controller: :users

    resources :languages, only: [:show] do
      scope module: :languages do
        resources :modules, only: []
      end
    end
  end
end
