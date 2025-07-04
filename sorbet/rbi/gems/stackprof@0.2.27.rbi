# typed: true

# DO NOT EDIT MANUALLY
# This is an autogenerated file for types exported from the `stackprof` gem.
# Please instead update this file by running `bin/tapioca gem stackprof`.


# source://stackprof//lib/stackprof.rb#20
module StackProf
  class << self
    # source://stackprof//lib/stackprof.rb#4
    def results(*_arg0); end

    # source://stackprof//lib/stackprof.rb#4
    def run(*_arg0); end

    # @return [Boolean]
    #
    # source://stackprof//lib/stackprof.rb#4
    def running?; end

    # source://stackprof//lib/stackprof.rb#4
    def sample; end

    # source://stackprof//lib/stackprof.rb#4
    def start(*_arg0); end

    # source://stackprof//lib/stackprof.rb#4
    def stop; end

    # source://stackprof//lib/stackprof.rb#4
    def use_postponed_job!; end
  end
end

# source://stackprof//lib/stackprof/middleware.rb#4
class StackProf::Middleware
  # @return [Middleware] a new instance of Middleware
  #
  # source://stackprof//lib/stackprof/middleware.rb#5
  def initialize(app, options = T.unsafe(nil)); end

  # source://stackprof//lib/stackprof/middleware.rb#20
  def call(env); end

  class << self
    # Returns the value of attribute enabled.
    #
    # source://stackprof//lib/stackprof/middleware.rb#40
    def enabled; end

    # Sets the attribute enabled
    #
    # @param value the value to set the attribute enabled to.
    #
    # source://stackprof//lib/stackprof/middleware.rb#40
    def enabled=(_arg0); end

    # @return [Boolean]
    #
    # source://stackprof//lib/stackprof/middleware.rb#42
    def enabled?(env); end

    # Returns the value of attribute interval.
    #
    # source://stackprof//lib/stackprof/middleware.rb#40
    def interval; end

    # Sets the attribute interval
    #
    # @param value the value to set the attribute interval to.
    #
    # source://stackprof//lib/stackprof/middleware.rb#40
    def interval=(_arg0); end

    # Returns the value of attribute metadata.
    #
    # source://stackprof//lib/stackprof/middleware.rb#40
    def metadata; end

    # Sets the attribute metadata
    #
    # @param value the value to set the attribute metadata to.
    #
    # source://stackprof//lib/stackprof/middleware.rb#40
    def metadata=(_arg0); end

    # Returns the value of attribute mode.
    #
    # source://stackprof//lib/stackprof/middleware.rb#40
    def mode; end

    # Sets the attribute mode
    #
    # @param value the value to set the attribute mode to.
    #
    # source://stackprof//lib/stackprof/middleware.rb#40
    def mode=(_arg0); end

    # Returns the value of attribute path.
    #
    # source://stackprof//lib/stackprof/middleware.rb#40
    def path; end

    # Sets the attribute path
    #
    # @param value the value to set the attribute path to.
    #
    # source://stackprof//lib/stackprof/middleware.rb#40
    def path=(_arg0); end

    # Returns the value of attribute raw.
    #
    # source://stackprof//lib/stackprof/middleware.rb#40
    def raw; end

    # Sets the attribute raw
    #
    # @param value the value to set the attribute raw to.
    #
    # source://stackprof//lib/stackprof/middleware.rb#40
    def raw=(_arg0); end

    # source://stackprof//lib/stackprof/middleware.rb#50
    def save; end
  end
end

# source://stackprof//lib/stackprof/report.rb#8
class StackProf::Report
  # @return [Report] a new instance of Report
  #
  # source://stackprof//lib/stackprof/report.rb#42
  def initialize(data); end

  # @raise [ArgumentError]
  #
  # source://stackprof//lib/stackprof/report.rb#618
  def +(other); end

  # source://stackprof//lib/stackprof/report.rb#92
  def add_lines(a, b); end

  # source://stackprof//lib/stackprof/report.rb#212
  def convert_to_d3_flame_graph_format(name, stacks, depth); end

  # Returns the value of attribute data.
  #
  # source://stackprof//lib/stackprof/report.rb#45
  def data; end

  # source://stackprof//lib/stackprof/report.rb#80
  def files; end

  # source://stackprof//lib/stackprof/report.rb#205
  def flamegraph_row(f, x, y, weight, addr); end

  # source://stackprof//lib/stackprof/report.rb#187
  def flamegraph_stacks(raw); end

  # source://stackprof//lib/stackprof/report.rb#47
  def frames(sort_by_total = T.unsafe(nil)); end

  # source://stackprof//lib/stackprof/report.rb#76
  def max_samples; end

  # source://stackprof//lib/stackprof/report.rb#68
  def modeline; end

  # source://stackprof//lib/stackprof/report.rb#52
  def normalized_frames; end

  # source://stackprof//lib/stackprof/report.rb#72
  def overall_samples; end

  # source://stackprof//lib/stackprof/report.rb#128
  def print_alphabetical_flamegraph(f = T.unsafe(nil), skip_common = T.unsafe(nil)); end

  # source://stackprof//lib/stackprof/report.rb#495
  def print_callgrind(f = T.unsafe(nil)); end

  # source://stackprof//lib/stackprof/report.rb#242
  def print_d3_flamegraph(f = T.unsafe(nil), skip_common = T.unsafe(nil)); end

  # source://stackprof//lib/stackprof/report.rb#99
  def print_debug; end

  # source://stackprof//lib/stackprof/report.rb#103
  def print_dump(f = T.unsafe(nil)); end

  # source://stackprof//lib/stackprof/report.rb#610
  def print_file(filter, f = T.unsafe(nil)); end

  # source://stackprof//lib/stackprof/report.rb#600
  def print_files(sort_by_total = T.unsafe(nil), limit = T.unsafe(nil), f = T.unsafe(nil)); end

  # source://stackprof//lib/stackprof/report.rb#132
  def print_flamegraph(f, skip_common, alphabetical = T.unsafe(nil)); end

  # source://stackprof//lib/stackprof/report.rb#414
  def print_graphviz(options = T.unsafe(nil), f = T.unsafe(nil)); end

  # source://stackprof//lib/stackprof/report.rb#107
  def print_json(f = T.unsafe(nil)); end

  # source://stackprof//lib/stackprof/report.rb#526
  def print_method(name, f = T.unsafe(nil)); end

  # source://stackprof//lib/stackprof/report.rb#112
  def print_stackcollapse; end

  # source://stackprof//lib/stackprof/report.rb#476
  def print_text(sort_by_total = T.unsafe(nil), limit = T.unsafe(nil), select_files = T.unsafe(nil), reject_files = T.unsafe(nil), select_names = T.unsafe(nil), reject_names = T.unsafe(nil), f = T.unsafe(nil)); end

  # source://stackprof//lib/stackprof/report.rb#124
  def print_timeline_flamegraph(f = T.unsafe(nil), skip_common = T.unsafe(nil)); end

  # source://stackprof//lib/stackprof/report.rb#64
  def version; end

  # Walk up and down the stack from a given starting point (name).  Loops
  # until `:exit` is selected
  #
  # source://stackprof//lib/stackprof/report.rb#561
  def walk_method(name); end

  private

  # source://stackprof//lib/stackprof/report.rb#669
  def callers_for(addr); end

  # source://stackprof//lib/stackprof/report.rb#665
  def root_frames; end

  # source://stackprof//lib/stackprof/report.rb#674
  def source_display(f, file, lines, range = T.unsafe(nil)); end

  class << self
    # source://stackprof//lib/stackprof/report.rb#12
    def from_file(file); end

    # source://stackprof//lib/stackprof/report.rb#20
    def from_json(json); end

    # source://stackprof//lib/stackprof/report.rb#24
    def parse_json(json); end
  end
end

# source://stackprof//lib/stackprof/report.rb#9
StackProf::Report::MARSHAL_SIGNATURE = T.let(T.unsafe(nil), String)

# source://stackprof//lib/stackprof.rb#21
StackProf::VERSION = T.let(T.unsafe(nil), String)
