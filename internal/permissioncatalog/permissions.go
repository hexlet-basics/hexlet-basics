// Package permissioncatalog is the catalog of Permissions an admin operation can
// require (ADR-0020). Each operation declares its Permission and whether it
// reads or writes with `x-permission` in the TypeSpec; this package reads those
// declarations out of the generated OpenAPI document, so the catalog is never
// kept by hand and cannot drift from the contract.
package permissioncatalog

import (
	"fmt"
	"slices"
	"strings"

	"gopkg.in/yaml.v3"

	apispec "hexletbasics/api-spec"
)

// Kind says whether an operation only reads or can change data. Over MCP it
// picks the `mcp.read` or `mcp.write` Permission an operation also needs.
type Kind string

const (
	Read  Kind = "read"
	Write Kind = "write"
)

// Access is the requirement an admin operation declares in place of a
// Permission, with `x-access` in the TypeSpec. Every operation that used to be
// administrators-only now declares a Permission, so what remains is the upload.
type Access string

const (
	// AnyStaff marks an operation any staff member may perform, whatever their
	// role holds: the upload every form uses.
	AnyStaff Access = "any-staff"
)

// Requirement is what one operation declares it needs.
type Requirement struct {
	Permission string
	Kind       Kind
}

// Permission is one entry of the catalog. A Permission has one Kind: it is
// either the right to read or the right to change, so a role can be given one
// without the other.
type Permission struct {
	Name string
	Kind Kind
	// Capability marks a Permission listed under `x-capabilities`, which the
	// administrator flag does not imply.
	Capability bool
}

// Declarations is what the contract declares: the requirement of each
// operation, keyed by its ogen name (AdminListBlogPosts), the operations that
// declare an Access instead, the capability Permissions, and the catalog.
type Declarations struct {
	operations   map[string]Requirement
	access       map[string]Access
	capabilities map[string]Kind
	adminWrite   map[string]bool
	standalone   map[string]bool
	personalData map[string]string
	required     map[string]string
	catalog      []Permission
}

// declared is read once from the embedded contract. A contract that cannot be
// parsed is a build defect, not a runtime condition, and TestTheCatalog...
// covers it, so a failure here panics at start-up rather than silently opening
// or closing every operation.
var declared = mustParse(apispec.OpenAPI)

func mustParse(doc []byte) *Declarations {
	d, err := Parse(doc)
	if err != nil {
		panic(fmt.Sprintf("permissions: read the contract: %v", err))
	}
	return d
}

// Parse reads the `x-permission` declarations out of an OpenAPI document. It is
// exported so tests can feed it malformed documents.
func Parse(doc []byte) (*Declarations, error) {
	var spec struct {
		// Capabilities is the root-level `x-capabilities` list: Permissions
		// nothing in the contract requires yet (mcp.read, audit.read).
		Capabilities []struct {
			Name string `yaml:"name"`
			Kind Kind   `yaml:"kind"`
		} `yaml:"x-capabilities"`
		Paths map[string]map[string]struct {
			OperationID  string                `yaml:"operationId"`
			Permission   *permissionDecl       `yaml:"x-permission"`
			Access       Access                `yaml:"x-access"`
			UnitOfWork   string                `yaml:"x-unit-of-work"`
			PersonalData string                `yaml:"x-personal-data"`
			Requires     string                `yaml:"x-requires-capability"`
			Security     []map[string][]string `yaml:"security"`
		} `yaml:"paths"`
	}
	if err := yaml.Unmarshal(doc, &spec); err != nil {
		return nil, fmt.Errorf("parse OpenAPI: %w", err)
	}

	d := &Declarations{
		operations:   map[string]Requirement{},
		access:       map[string]Access{},
		capabilities: map[string]Kind{},
		adminWrite:   map[string]bool{},
		standalone:   map[string]bool{},
		personalData: map[string]string{},
		required:     map[string]string{},
	}
	kinds := map[string]Kind{}
	for _, c := range spec.Capabilities {
		if c.Name == "" || (c.Kind != Read && c.Kind != Write) {
			return nil, fmt.Errorf("x-capabilities: %q needs a name and a kind of read or write", c.Name)
		}
		if _, dup := kinds[c.Name]; dup {
			return nil, fmt.Errorf("x-capabilities: %s is listed twice", c.Name)
		}
		kinds[c.Name] = c.Kind
		d.capabilities[c.Name] = c.Kind
	}
	for path, methods := range spec.Paths {
		for method, op := range methods {
			where := fmt.Sprintf("%s %s", strings.ToUpper(method), path)
			if op.Access != "" {
				if op.Access != AnyStaff {
					return nil, fmt.Errorf("%s: x-access %q is not %s", where, op.Access, AnyStaff)
				}
				if op.Permission != nil {
					return nil, fmt.Errorf("%s: declares both x-permission and x-access", where)
				}
				if op.OperationID == "" {
					return nil, fmt.Errorf("%s: x-access needs an operationId", where)
				}
				d.access[operationName(op.OperationID)] = op.Access
			}
			if op.UnitOfWork != "" {
				if op.UnitOfWork != unitStandalone {
					return nil, fmt.Errorf("%s: x-unit-of-work %q is not %s", where, op.UnitOfWork, unitStandalone)
				}
				if op.OperationID == "" {
					return nil, fmt.Errorf("%s: x-unit-of-work needs an operationId", where)
				}
				d.standalone[operationName(op.OperationID)] = true
			}
			if op.Requires != "" {
				if op.OperationID == "" {
					return nil, fmt.Errorf("%s: x-requires-capability needs an operationId", where)
				}
				if _, ok := d.capabilities[op.Requires]; !ok {
					return nil, fmt.Errorf("%s: x-requires-capability %q is not listed under x-capabilities", where, op.Requires)
				}
				d.required[operationName(op.OperationID)] = op.Requires
			}
			if op.PersonalData != "" {
				if op.OperationID == "" {
					return nil, fmt.Errorf("%s: x-personal-data needs an operationId", where)
				}
				if op.Permission == nil || op.Permission.Kind != Read {
					return nil, fmt.Errorf("%s: x-personal-data belongs on an operation that declares a read Permission", where)
				}
				d.personalData[operationName(op.OperationID)] = op.PersonalData
			}
			if op.OperationID != "" && isAdminWrite(method, op.Security, op.Permission != nil, op.Access != "", permissionKind(op.Permission)) {
				d.adminWrite[operationName(op.OperationID)] = true
			}
			if op.Permission == nil {
				continue
			}
			if op.OperationID == "" || op.Permission.Name == "" {
				return nil, fmt.Errorf("%s: x-permission needs an operationId and a name", where)
			}
			if op.Permission.Kind != Read && op.Permission.Kind != Write {
				return nil, fmt.Errorf("%s: x-permission kind %q is neither read nor write", where, op.Permission.Kind)
			}
			if known, ok := kinds[op.Permission.Name]; ok && known != op.Permission.Kind {
				return nil, fmt.Errorf("%s: %s is declared both %s and %s", where, op.Permission.Name, known, op.Permission.Kind)
			}
			kinds[op.Permission.Name] = op.Permission.Kind
			d.operations[operationName(op.OperationID)] = Requirement{Permission: op.Permission.Name, Kind: op.Permission.Kind}
		}
	}

	for name, kind := range kinds {
		_, capability := d.capabilities[name]
		d.catalog = append(d.catalog, Permission{Name: name, Kind: kind, Capability: capability})
	}
	slices.SortFunc(d.catalog, func(a, b Permission) int { return strings.Compare(a.Name, b.Name) })
	return d, nil
}

// unitStandalone is the only value of `x-unit-of-work`: the operation calls a
// service outside the database, so the audit middleware must not hold a
// transaction open around it.
const unitStandalone = "standalone"

// permissionDecl is the `x-permission` value of one operation.
type permissionDecl struct {
	Name string `yaml:"name"`
	Kind Kind   `yaml:"kind"`
}

func permissionKind(p *permissionDecl) Kind {
	if p == nil {
		return ""
	}
	return p.Kind
}

// isAdminWrite is the one rule that says an operation changes data on the admin
// surface, so the audit trail needs no list of its own. An operation is admin
// when it requires the AdminSession security scheme or declares a Permission or
// an Access (both only exist on admin operations). It is a write when its
// declared Permission says so; one that declares only an Access (or nothing yet)
// is judged by its HTTP method, so an operation marked administrators-only is
// audited before #843 gives it a Permission, and stays audited after.
func isAdminWrite(method string, security []map[string][]string, hasPermission, hasAccess bool, kind Kind) bool {
	admin := hasPermission || hasAccess
	for _, requirement := range security {
		if _, ok := requirement["AdminSession"]; ok {
			admin = true
		}
	}
	if !admin {
		return false
	}
	if kind != "" {
		return kind == Write
	}
	switch strings.ToUpper(method) {
	case "GET", "HEAD", "OPTIONS":
		return false
	}
	return true
}

// operationName turns an OpenAPI operationId (adminListBanners) into the ogen
// name (AdminListBanners) the security handler receives.
func operationName(operationID string) string {
	return strings.ToUpper(operationID[:1]) + operationID[1:]
}

// For reports what the operation (an ogen operation name) declares it needs.
// ok is false for an operation that declares no Permission; AccessOf says
// whether it declares an Access instead.
func For(operation string) (Requirement, bool) {
	r, ok := declared.operations[operation]
	return r, ok
}

// AccessOf reports the Access an operation declares instead of a Permission, and
// the empty Access for one that declares none.
func AccessOf(operation string) Access {
	return declared.access[operation]
}

// IsClassified reports whether an operation says who may use it: it declares a
// Permission or an Access. An admin operation that does neither is closed to
// everyone but administrators by default, and a contract test fails on it so it
// is never closed by accident.
func IsClassified(operation string) bool {
	_, hasPermission := declared.operations[operation]
	return hasPermission || declared.access[operation] != ""
}

// IsAdminWrite reports whether an operation changes data on the admin surface:
// the operations whose every call must leave an Audit Record (ADR-0020). It is
// read from the contract's security requirement, Permission kind and HTTP
// method, never from a Go list.
func IsAdminWrite(operation string) bool {
	return declared.IsAdminWrite(operation)
}

// AdminWrites lists every operation IsAdminWrite accepts, sorted, so a test can
// walk the whole admin write surface.
func AdminWrites() []string {
	return declared.AdminWrites()
}

// Standalone reports whether an operation declares `x-unit-of-work: standalone`:
// it calls a service outside the database (blob storage, a model, a mailer), so
// it must not run inside a transaction held open for the call.
func Standalone(operation string) bool {
	return declared.Standalone(operation)
}

// IsAdminWrite is the package-level function over these declarations.
func (d *Declarations) IsAdminWrite(operation string) bool { return d.adminWrite[operation] }

// AdminWrites is the package-level function over these declarations.
func (d *Declarations) AdminWrites() []string {
	names := make([]string, 0, len(d.adminWrite))
	for name := range d.adminWrite {
		names = append(names, name)
	}
	slices.Sort(names)
	return names
}

// Standalone is the package-level function over these declarations.
func (d *Declarations) Standalone(operation string) bool { return d.standalone[operation] }

// PersonalDataEntity reports whether an operation reads personal data, as the
// contract declares with `x-personal-data`, and the ent type it reads (User,
// Lead). Every such read leaves an Audit Record (ADR-0020); a read of the
// catalogue or of content declares nothing and leaves none.
func PersonalDataEntity(operation string) (entity string, ok bool) {
	return declared.PersonalDataEntity(operation)
}

// PersonalDataReads lists every operation that declares `x-personal-data`,
// sorted, so a test can walk the whole set.
func PersonalDataReads() []string {
	return declared.PersonalDataReads()
}

// PersonalDataEntity is the package-level function over these declarations.
func (d *Declarations) PersonalDataEntity(operation string) (string, bool) {
	entity, ok := d.personalData[operation]
	return entity, ok
}

// PersonalDataReads is the package-level function over these declarations.
func (d *Declarations) PersonalDataReads() []string {
	names := make([]string, 0, len(d.personalData))
	for name := range d.personalData {
		names = append(names, name)
	}
	slices.Sort(names)
	return names
}

// Catalog lists every declared Permission once, ordered by name so the roles
// screen and the tests see a stable list.
func Catalog() []Permission {
	return slices.Clone(declared.catalog)
}

// IsCapability reports whether name is a capability Permission, one the
// contract lists under `x-capabilities`. The administrator flag does not imply
// these: a role must be given them explicitly, administrators' included.
func IsCapability(name string) bool {
	_, ok := declared.capabilities[name]
	return ok
}

// Known reports whether name is in the catalog. Grants are validated against it
// so a role cannot be given a Permission nothing requires.
func Known(name string) bool {
	return slices.ContainsFunc(declared.catalog, func(p Permission) bool { return p.Name == name })
}

// MCPCapabilities are the capability Permissions an operation needs on top of
// its own when it is called over MCP (ADR-0020), all of them: mcp.write for an
// operation that changes data on the admin surface, mcp.read for the rest, and
// the capability the contract declares with `x-requires-capability` (mcp.manage
// for every operation on users, staff members, Staff Roles and their grants,
// whatever its kind). The kind comes from the contract (the declared
// Permission's, else the one IsAdminWrite derives), so a new operation is
// classified by declaring it, never by a list in Go.
func MCPCapabilities(operation string) []string {
	kind := "mcp.read"
	if IsAdminWrite(operation) {
		kind = "mcp.write"
	}
	if extra, ok := declared.required[operation]; ok {
		return []string{kind, extra}
	}
	return []string{kind}
}

// RequiredCapabilityOperations lists the operations that declare
// `x-requires-capability`, sorted, so a test can walk them.
func RequiredCapabilityOperations() []string {
	names := make([]string, 0, len(declared.required))
	for name := range declared.required {
		names = append(names, name)
	}
	slices.Sort(names)
	return names
}
