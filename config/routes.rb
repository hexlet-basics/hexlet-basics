# frozen_string_literal: true

Rails.application.routes.draw do
  scope module: :web do
    root 'home#index'

    resources :languages, only: [:show] do
      scope module: :languages do
        resources :modules
      end
    end
  end
end
