@0xda7a3c24c13f41ac;  # Unique file ID, you should generate a new one

struct Artifact {
  header @0 : Header;
  versions @1 : List(PayloadVersion);

  struct Header {
    id @0 : Text;
    schemaVersion @1 : UInt32;
    role @2 : Text;
    type @3 : Text;
    tags @4 : List(Text);
    scopes @5 : List(Text);
    status @6 : Text;
    auditTrail @7 : List(AuditEntry);
    securityInfo @8 : SecurityMetadata;
  }

  struct SecurityMetadata {
    encryptionMethod @0 : Text;
    allowedUsers @1 : List(Text);
    allowedGroups @2 : List(Text);
    publicKey @3 : Text;
    expirationDate @4 : Int64;  # Unix timestamp
  }

  struct AuditEntry {
    timestamp @0 : Int64;  # Unix timestamp
    action @1 : Text;
    performedBy @2 : Text;
    systemId @3 : Text;
    additionalInfo @4 : Text;  # JSON-encoded string for flexibility
    digitalSignature @5 : Text;
  }

  struct PayloadVersion {
    versionNumber @0 : UInt32;
    timestamp @1 : Int64;  # Unix timestamp
    modifiedBy @2 : Text;
    changeDescription @3 : Text;
    data @4 : List(PayloadData);
    integrityInfo @5 : IntegrityCheck;
  }

  struct PayloadData {
    encryptedData @0 : Data;
    encryptionKeyId @1 : Text;
  }

  struct IntegrityCheck {
    hashMethod @0 : Text;
    dataHash @1 : Text;
    digitalSignature @2 : Text;
  }
}
