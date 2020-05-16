# frozen_string_literal: true

Rails.application.routes.draw do
  scope module: :web do
    root 'home#index'

    resources :sessions, only: [:new, :create, :destroy]
    resources :registrations, only: [:new, :create], controller: :users

    resources :languages, only: [:show] do
      scope module: :languages do
        resources :modules, only: []
      end
    end
  end
end
