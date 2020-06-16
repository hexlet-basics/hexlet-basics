# frozen_string_literal: true

Rails.application.routes.draw do
  require 'sidekiq/web'
  mount Sidekiq::Web => '/sidekiq', constraints: AdminConstraint.new

  scope module: :web do
    root 'home#index'

    resource :session, only: %i[new create destroy]
    resources :users, only: %i[new create]

    resources :languages, only: [:show] do
      scope module: :languages do
        resources :modules, only: [] do
          scope module: :modules do
            resources :lessons, only: [:show]
          end
        end
      end
    end

    namespace :admin do
      root 'home#index'

      resources :uploads, only: %i[index new create]
      resources :languages, only: %i[index new create] do
        scope module: :languages do
          resources :uploads, only: [:create]
        end
      end
    end
  end
end
