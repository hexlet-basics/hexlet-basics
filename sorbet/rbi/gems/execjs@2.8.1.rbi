# typed: true

# DO NOT EDIT MANUALLY
# This is an autogenerated file for types exported from the `execjs` gem.
# Please instead update this file by running `bin/tapioca gem execjs`.

module ExecJS
  class << self
    # source://execjs//lib/execjs/module.rb#26
    def compile(source, options = T.unsafe(nil)); end

    # @return [Boolean]
    #
    # source://execjs//lib/execjs/module.rb#38
    def cygwin?; end

    # source://execjs//lib/execjs/module.rb#22
    def eval(source, options = T.unsafe(nil)); end

    # source://execjs//lib/execjs/module.rb#18
    def exec(source, options = T.unsafe(nil)); end

    # source://execjs//lib/execjs/module.rb#30
    def root; end

    # Returns the value of attribute runtime.
    #
    # source://execjs//lib/execjs/module.rb#11
    def runtime; end

    # @raise [RuntimeUnavailable]
    #
    # source://execjs//lib/execjs/module.rb#13
    def runtime=(runtime); end

    # source://execjs//lib/execjs/runtimes.rb#96
    def runtimes; end

    # @return [Boolean]
    #
    # source://execjs//lib/execjs/module.rb#34
    def windows?; end
  end
end

class ExecJS::DisabledRuntime < ::ExecJS::Runtime
  # @return [Boolean]
  #
  # source://execjs//lib/execjs/disabled_runtime.rb#25
  def available?; end

  # @raise [Error]
  #
  # source://execjs//lib/execjs/disabled_runtime.rb#17
  def compile(source, options = T.unsafe(nil)); end

  # @return [Boolean]
  #
  # source://execjs//lib/execjs/disabled_runtime.rb#21
  def deprecated?; end

  # @raise [Error]
  #
  # source://execjs//lib/execjs/disabled_runtime.rb#13
  def eval(source, options = T.unsafe(nil)); end

  # @raise [Error]
  #
  # source://execjs//lib/execjs/disabled_runtime.rb#9
  def exec(source, options = T.unsafe(nil)); end

  # source://execjs//lib/execjs/disabled_runtime.rb#5
  def name; end
end

class ExecJS::DuktapeRuntime < ::ExecJS::Runtime
  # @return [Boolean]
  #
  # source://execjs//lib/execjs/duktape_runtime.rb#62
  def available?; end

  # source://execjs//lib/execjs/duktape_runtime.rb#58
  def name; end
end

class ExecJS::DuktapeRuntime::Context < ::ExecJS::Runtime::Context
  # @return [Context] a new instance of Context
  #
  # source://execjs//lib/execjs/duktape_runtime.rb#7
  def initialize(runtime, source = T.unsafe(nil), options = T.unsafe(nil)); end

  # source://execjs//lib/execjs/duktape_runtime.rb#28
  def call(identifier, *args); end

  # source://execjs//lib/execjs/duktape_runtime.rb#21
  def eval(source, options = T.unsafe(nil)); end

  # source://execjs//lib/execjs/duktape_runtime.rb#14
  def exec(source, options = T.unsafe(nil)); end

  private

  # source://execjs//lib/execjs/duktape_runtime.rb#36
  def wrap_error(e); end
end

# Encodes strings as UTF-8
module ExecJS::Encoding
  # workaround for jruby bug http://jira.codehaus.org/browse/JRUBY-6588
  # workaround for rbx bug https://github.com/rubinius/rubinius/issues/1729
  #
  # source://execjs//lib/execjs/encoding.rb#21
  def encode(string); end
end

class ExecJS::Error < ::StandardError; end

class ExecJS::ExternalRuntime < ::ExecJS::Runtime
  # @return [ExternalRuntime] a new instance of ExternalRuntime
  #
  # source://execjs//lib/execjs/external_runtime.rb#92
  def initialize(options); end

  # @return [Boolean]
  #
  # source://execjs//lib/execjs/external_runtime.rb#109
  def available?; end

  # @return [Boolean]
  #
  # source://execjs//lib/execjs/external_runtime.rb#114
  def deprecated?; end

  # source://execjs//lib/execjs/external_runtime.rb#211
  def exec_runtime(filename); end

  # Returns the value of attribute name.
  #
  # source://execjs//lib/execjs/external_runtime.rb#90
  def name; end

  protected

  # source://execjs//lib/execjs/external_runtime.rb#159
  def encode_source(source); end

  # source://execjs//lib/execjs/external_runtime.rb#164
  def encode_unicode_codepoints(str); end

  # source://execjs//lib/execjs/external_runtime.rb#226
  def exec_runtime_error(output); end

  # source://execjs//lib/execjs/external_runtime.rb#145
  def generate_compile_method(path); end

  # source://execjs//lib/execjs/external_runtime.rb#155
  def json2_source; end

  # source://execjs//lib/execjs/external_runtime.rb#235
  def which(command); end

  private

  # source://execjs//lib/execjs/external_runtime.rb#119
  def binary; end

  # source://execjs//lib/execjs/external_runtime.rb#123
  def locate_executable(command); end
end

class ExecJS::ExternalRuntime::Context < ::ExecJS::Runtime::Context
  # @return [Context] a new instance of Context
  #
  # source://execjs//lib/execjs/external_runtime.rb#7
  def initialize(runtime, source = T.unsafe(nil), options = T.unsafe(nil)); end

  # source://execjs//lib/execjs/external_runtime.rb#45
  def call(identifier, *args); end

  # source://execjs//lib/execjs/external_runtime.rb#17
  def eval(source, options = T.unsafe(nil)); end

  # source://execjs//lib/execjs/external_runtime.rb#25
  def exec(source, options = T.unsafe(nil)); end

  protected

  # See Tempfile.create on Ruby 2.1
  #
  # source://execjs//lib/execjs/external_runtime.rb#51
  def create_tempfile(basename); end

  # source://execjs//lib/execjs/external_runtime.rb#67
  def extract_result(output, filename); end

  # source://execjs//lib/execjs/external_runtime.rb#60
  def write_to_tempfile(contents); end
end

class ExecJS::MiniRacerRuntime < ::ExecJS::Runtime
  # @return [Boolean]
  #
  # source://execjs//lib/execjs/mini_racer_runtime.rb#96
  def available?; end

  # source://execjs//lib/execjs/mini_racer_runtime.rb#92
  def name; end
end

class ExecJS::MiniRacerRuntime::Context < ::ExecJS::Runtime::Context
  # @return [Context] a new instance of Context
  #
  # source://execjs//lib/execjs/mini_racer_runtime.rb#6
  def initialize(runtime, source = T.unsafe(nil), options = T.unsafe(nil)); end

  # source://execjs//lib/execjs/mini_racer_runtime.rb#33
  def call(identifier, *args); end

  # source://execjs//lib/execjs/mini_racer_runtime.rb#23
  def eval(source, options = T.unsafe(nil)); end

  # source://execjs//lib/execjs/mini_racer_runtime.rb#15
  def exec(source, options = T.unsafe(nil)); end

  private

  # source://execjs//lib/execjs/mini_racer_runtime.rb#40
  def strip_functions!(value); end

  # source://execjs//lib/execjs/mini_racer_runtime.rb#65
  def translate; end
end

class ExecJS::ProgramError < ::ExecJS::Error; end

class ExecJS::RubyRhinoRuntime < ::ExecJS::Runtime
  # @return [Boolean]
  #
  # source://execjs//lib/execjs/ruby_rhino_runtime.rb#90
  def available?; end

  # source://execjs//lib/execjs/ruby_rhino_runtime.rb#86
  def name; end
end

class ExecJS::RubyRhinoRuntime::Context < ::ExecJS::Runtime::Context
  # @return [Context] a new instance of Context
  #
  # source://execjs//lib/execjs/ruby_rhino_runtime.rb#6
  def initialize(runtime, source = T.unsafe(nil), options = T.unsafe(nil)); end

  # source://execjs//lib/execjs/ruby_rhino_runtime.rb#34
  def call(properties, *args); end

  # source://execjs//lib/execjs/ruby_rhino_runtime.rb#24
  def eval(source, options = T.unsafe(nil)); end

  # source://execjs//lib/execjs/ruby_rhino_runtime.rb#16
  def exec(source, options = T.unsafe(nil)); end

  # source://execjs//lib/execjs/ruby_rhino_runtime.rb#40
  def unbox(value); end

  # source://execjs//lib/execjs/ruby_rhino_runtime.rb#61
  def wrap_error(e); end

  private

  # Disables bytecode compiling which limits you to 64K scripts
  #
  # source://execjs//lib/execjs/ruby_rhino_runtime.rb#77
  def fix_memory_limit!(context); end
end

# Abstract base class for runtimes
class ExecJS::Runtime
  # @raise [NotImplementedError]
  # @return [Boolean]
  #
  # source://execjs//lib/execjs/runtime.rb#80
  def available?; end

  # source://execjs//lib/execjs/runtime.rb#68
  def compile(source, options = T.unsafe(nil)); end

  # source://execjs//lib/execjs/runtime.rb#44
  def context_class; end

  # @return [Boolean]
  #
  # source://execjs//lib/execjs/runtime.rb#76
  def deprecated?; end

  # source://execjs//lib/execjs/runtime.rb#58
  def eval(source, options = T.unsafe(nil)); end

  # source://execjs//lib/execjs/runtime.rb#48
  def exec(source, options = T.unsafe(nil)); end

  # @raise [NotImplementedError]
  #
  # source://execjs//lib/execjs/runtime.rb#40
  def name; end
end

class ExecJS::Runtime::Context
  include ::ExecJS::Encoding

  # @return [Context] a new instance of Context
  #
  # source://execjs//lib/execjs/runtime.rb#9
  def initialize(runtime, source = T.unsafe(nil), options = T.unsafe(nil)); end

  # Evaluates +source+ as an expression (which should be of type
  # +function+), and calls the function with the given arguments.
  # The function will be evaluated with the global object as +this+.
  #
  #   context.call("function(a, b) { return a + b }", 1, 1) # => 2
  #   context.call("CoffeeScript.compile", "1 + 1")
  #
  # @raise [NotImplementedError]
  #
  # source://execjs//lib/execjs/runtime.rb#35
  def call(source, *args); end

  # Evaluates the +source+ as an expression and returns the result.
  #
  #   context.eval("1")        # => 1
  #   context.eval("return 1") # => Raises SyntaxError
  #
  # @raise [NotImplementedError]
  #
  # source://execjs//lib/execjs/runtime.rb#25
  def eval(source, options = T.unsafe(nil)); end

  # Evaluates the +source+ in the context of a function body and returns the
  # returned value.
  #
  #   context.exec("return 1") # => 1
  #   context.exec("1")        # => nil (nothing was returned)
  #
  # @raise [NotImplementedError]
  #
  # source://execjs//lib/execjs/runtime.rb#17
  def exec(source, options = T.unsafe(nil)); end
end

class ExecJS::RuntimeError < ::ExecJS::Error; end
class ExecJS::RuntimeUnavailable < ::ExecJS::RuntimeError; end

module ExecJS::Runtimes
  class << self
    # source://execjs//lib/execjs/runtimes.rb#56
    def autodetect; end

    # source://execjs//lib/execjs/runtimes.rb#62
    def best_available; end

    # source://execjs//lib/execjs/runtimes.rb#66
    def from_environment; end

    # source://execjs//lib/execjs/runtimes.rb#78
    def names; end

    # source://execjs//lib/execjs/runtimes.rb#82
    def runtimes; end
  end
end

# source://execjs//lib/execjs/runtimes.rb#10
ExecJS::Runtimes::Disabled = T.let(T.unsafe(nil), ExecJS::DisabledRuntime)

# source://execjs//lib/execjs/runtimes.rb#12
ExecJS::Runtimes::Duktape = T.let(T.unsafe(nil), ExecJS::DuktapeRuntime)

# source://execjs//lib/execjs/runtimes.rb#41
ExecJS::Runtimes::JScript = T.let(T.unsafe(nil), ExecJS::ExternalRuntime)

# source://execjs//lib/execjs/runtimes.rb#25
ExecJS::Runtimes::JavaScriptCore = T.let(T.unsafe(nil), ExecJS::ExternalRuntime)

# source://execjs//lib/execjs/runtimes.rb#16
ExecJS::Runtimes::MiniRacer = T.let(T.unsafe(nil), ExecJS::MiniRacerRuntime)

# source://execjs//lib/execjs/runtimes.rb#18
ExecJS::Runtimes::Node = T.let(T.unsafe(nil), ExecJS::ExternalRuntime)

# source://execjs//lib/execjs/runtimes.rb#14
ExecJS::Runtimes::RubyRhino = T.let(T.unsafe(nil), ExecJS::RubyRhinoRuntime)

# source://execjs//lib/execjs/runtimes.rb#34
ExecJS::Runtimes::SpiderMonkey = T.let(T.unsafe(nil), ExecJS::ExternalRuntime)

# source://execjs//lib/execjs/runtimes.rb#34
ExecJS::Runtimes::Spidermonkey = T.let(T.unsafe(nil), ExecJS::ExternalRuntime)

# source://execjs//lib/execjs/runtimes.rb#48
ExecJS::Runtimes::V8 = T.let(T.unsafe(nil), ExecJS::ExternalRuntime)

# source://execjs//lib/execjs/version.rb#2
ExecJS::VERSION = T.let(T.unsafe(nil), String)
