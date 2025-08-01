module TestHelpers
  def self.read_fixture_file(name)
    content = Rails.root.join("test/fixtures/files", name)
    content.each_line.map { |line| "  #{line}" }.join
  end
end
