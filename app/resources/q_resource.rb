class QResource
  include Alba::Resource
  include Typelizer::DSL

  attributes :sf, :so, :fields
  typelize sf: :string, so: "SortOrder", fields: "Record<string, string | number | string[] | undefined>"
end
