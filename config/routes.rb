Rails.application.routes.draw do

  scope module: :web do
    root 'home#index'
  end
end
