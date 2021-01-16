# frozen_string_literal: true

class Open3Stub
  def self.capture2(_)
    output = 'Test passed!'
    process_status = OpenStruct.new(exitstatus: 0)

    [output, process_status]
  end
end
