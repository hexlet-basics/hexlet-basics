# frozen_string_literal: true

class ApplicationMailer < ActionMailer::Base
  default from: configus.mailer.from
  layout 'mailer'
end
