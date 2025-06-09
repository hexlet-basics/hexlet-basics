class RuSuffixConstraint
  def self.matches?(request)
    request.params[:suffix] == "ru"
  end
end
