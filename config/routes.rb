# frozen_string_literal: true

Rails.application.routes.draw do
  scope module: :web do
    root 'home#index'
  end
end
