package vmci

// VersionNumber is vSockets version.
//
// The version is a 32-bit unsigned integer that consist of three components:
// the epoch, the major version, and the minor version.
//
// Use Epoch, Major and Minor methods to extract the components.
type VersionNumber uint32

// UInt32 converts value to uint32 value.
//
// Alias to uint32(v).
func (v VersionNumber) UInt32() uint32 {
	return uint32(v)
}

// Epoch returns the epoch (first) component of the vSockets version.
//
// A single byte representing the epoch component of the vSockets version.
func (v VersionNumber) Epoch() uint32 {
	return ((v.UInt32()) & 0xFF000000) >> 24
}

// Major returns the major (second) component of the vSockets version.
//
// A single byte representing the major component of the vSockets version.
// Typically changes for every major release of a product
func (v VersionNumber) Major() uint32 {
	return ((v.UInt32()) & 0x00FF0000) >> 16
}

// Minor returns the minor (third) component of the vSockets version.
//
// Two bytes representing the minor component of the vSockets version.
func (v VersionNumber) Minor() uint32 {
	return v.UInt32() & 0x0000FFFF
}
