# typed: true

# DO NOT EDIT MANUALLY
# This is an autogenerated file for types exported from the `tins` gem.
# Please instead update this file by running `bin/tapioca gem tins`.

module Tins; end
module Tins::Blank; end

module Tins::Blank::Array
  class << self
    # @private
    #
    # source://tins//lib/tins/xt/blank.rb#32
    def included(modul); end
  end
end

module Tins::Blank::FalseClass
  # @return [Boolean]
  #
  # source://tins//lib/tins/xt/blank.rb#20
  def blank?; end
end

module Tins::Blank::Hash
  class << self
    # @private
    #
    # source://tins//lib/tins/xt/blank.rb#40
    def included(modul); end
  end
end

module Tins::Blank::NilClass
  # @return [Boolean]
  #
  # source://tins//lib/tins/xt/blank.rb#14
  def blank?; end
end

module Tins::Blank::Numeric
  # @return [Boolean]
  #
  # source://tins//lib/tins/xt/blank.rb#54
  def blank?; end
end

module Tins::Blank::Object
  # @return [Boolean]
  #
  # source://tins//lib/tins/xt/blank.rb#4
  def blank?; end

  # @return [Boolean]
  #
  # source://tins//lib/tins/xt/blank.rb#8
  def present?; end
end

module Tins::Blank::String
  # @return [Boolean]
  #
  # source://tins//lib/tins/xt/blank.rb#48
  def blank?; end
end

module Tins::Blank::TrueClass
  # @return [Boolean]
  #
  # source://tins//lib/tins/xt/blank.rb#26
  def blank?; end
end

module Tins::Full
  # @return [Boolean]
  #
  # source://tins//lib/tins/xt/full.rb#29
  def all_full?; end

  # Returns the object if it isn't blank (as in Object#blank?), otherwise it
  # returns nil. If a block was given as an argument and the object isn't
  # blank, the block is executed with the object as its first argument. If an
  # argument +dispatch+ was given and the object wasn't blank the method
  # given by dispatch is called on the object. This is the same as
  # foo.full?(&:bar) in the previous block form.
  #
  # @return [Boolean]
  #
  # source://tins//lib/tins/xt/full.rb#11
  def full?(dispatch = T.unsafe(nil), *args); end
end

module Tins::Terminal
  private

  # source://tins//lib/tins/terminal.rb#38
  def cols; end

  # source://tins//lib/tins/terminal.rb#33
  def columns; end

  # source://tins//lib/tins/terminal.rb#29
  def lines; end

  # source://tins//lib/tins/terminal.rb#24
  def rows; end

  # source://tins//lib/tins/terminal.rb#11
  def winsize; end

  class << self
    # source://tins//lib/tins/terminal.rb#38
    def cols; end

    # source://tins//lib/tins/terminal.rb#33
    def columns; end

    # source://tins//lib/tins/terminal.rb#29
    def lines; end

    # source://tins//lib/tins/terminal.rb#24
    def rows; end

    # source://tins//lib/tins/terminal.rb#11
    def winsize; end
  end
end
