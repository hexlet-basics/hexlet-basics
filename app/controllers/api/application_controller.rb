# typed: strict

class Api::ApplicationController < ApplicationController
  respond_to :json, :xml

  rescue_from ApplicationParamsStruct::InvalidParams do
    T.bind(self, Api::ApplicationController)
    head :unprocessable_entity
  end
end
