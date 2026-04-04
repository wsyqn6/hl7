package hl7

import "time"

type SegmentDefinition struct {
	Name          string
	Fields        []FieldDefinition
	IsRequired    bool
	MinOccurrence int
	MaxOccurrence int
}

type FieldDefinition struct {
	Index      int
	Name       string
	DataType   string
	IsRequired bool
	MaxLength  int
	Table      string
	Components []FieldDefinition
}

var SegmentDefinitions = map[string]SegmentDefinition{
	"MSH": {
		Name:          "MSH",
		IsRequired:    true,
		MinOccurrence: 1,
		MaxOccurrence: 1,
		Fields: []FieldDefinition{
			{Index: 1, Name: "FieldSeparator", DataType: "ST", IsRequired: true},
			{Index: 2, Name: "EncodingCharacters", DataType: "ST", IsRequired: true},
			{Index: 3, Name: "SendingApplication", DataType: "HD"},
			{Index: 4, Name: "SendingFacility", DataType: "HD"},
			{Index: 5, Name: "ReceivingApplication", DataType: "HD"},
			{Index: 6, Name: "ReceivingFacility", DataType: "HD"},
			{Index: 7, Name: "DateTimeOfMessage", DataType: "TS"},
			{Index: 8, Name: "Security", DataType: "ST"},
			{Index: 9, Name: "MessageType", DataType: "MSG", IsRequired: true},
			{Index: 10, Name: "MessageControlID", DataType: "ST"},
			{Index: 11, Name: "ProcessingID", DataType: "PT"},
			{Index: 12, Name: "VersionID", DataType: "VID"},
			{Index: 13, Name: "SequenceNumber", DataType: "NM"},
			{Index: 14, Name: "ContinuationPointer", DataType: "ST"},
			{Index: 15, Name: "AcceptAcknowledgmentType", DataType: "ID", Table: "HL70015"},
			{Index: 16, Name: "ApplicationAcknowledgmentType", DataType: "ID", Table: "HL70015"},
			{Index: 17, Name: "CountryCode", DataType: "ID", Table: "HL70039"},
		},
	},
	"PID": {
		Name:          "PID",
		IsRequired:    true,
		MinOccurrence: 1,
		MaxOccurrence: 1,
		Fields: []FieldDefinition{
			{Index: 1, Name: "SetID", DataType: "SI"},
			{Index: 2, Name: "PatientID", DataType: "CX"},
			{Index: 3, Name: "PatientIdentifierList", DataType: "CX", IsRequired: true, MaxLength: 20},
			{Index: 4, Name: "AlternatePatientID", DataType: "CX"},
			{Index: 5, Name: "PatientName", DataType: "XPN", IsRequired: true, Components: []FieldDefinition{
				{Index: 1, Name: "FamilyName", DataType: "FN"},
				{Index: 2, Name: "GivenName", DataType: "ST"},
				{Index: 3, Name: "SecondAndFurtherGivenNames", DataType: "ST"},
				{Index: 4, Name: "Suffix", DataType: "ST"},
				{Index: 5, Name: "Prefix", DataType: "ST"},
				{Index: 6, Name: "Degree", DataType: "ST"},
			}},
			{Index: 6, Name: "MothersMaidenName", DataType: "XPN"},
			{Index: 7, Name: "DateOfBirth", DataType: "TS"},
			{Index: 8, Name: "Sex", DataType: "ID", IsRequired: true, Table: "HL70001"},
			{Index: 9, Name: "PatientAlias", DataType: "XPN"},
			{Index: 10, Name: "Race", DataType: "CE", Table: "HL70005"},
			{Index: 11, Name: "PatientAddress", DataType: "XAD", Components: []FieldDefinition{
				{Index: 1, Name: "StreetAddress", DataType: "ST"},
				{Index: 2, Name: "OtherDesignation", DataType: "ST"},
				{Index: 3, Name: "City", DataType: "ST"},
				{Index: 4, Name: "StateOrProvince", DataType: "ST"},
				{Index: 5, Name: "ZipOrPostalCode", DataType: "ST"},
				{Index: 6, Name: "Country", DataType: "ID", Table: "HL70003"},
				{Index: 7, Name: "AddressType", DataType: "ID", Table: "HL70190"},
			}},
			{Index: 12, Name: "CountyCode", DataType: "ID", Table: "HL70179"},
			{Index: 13, Name: "PhoneHome", DataType: "XTN"},
			{Index: 14, Name: "PhoneBusiness", DataType: "XTN"},
			{Index: 15, Name: "PrimaryLanguage", DataType: "CE", Table: "HL70296"},
			{Index: 16, Name: "MaritalStatus", DataType: "CE", Table: "HL70002"},
			{Index: 17, Name: "Religion", DataType: "CE", Table: "HL70006"},
			{Index: 18, Name: "PatientAccountNumber", DataType: "CX"},
			{Index: 19, Name: "SSN", DataType: "ST"},
			{Index: 20, Name: "DriversLicense", DataType: "DLN"},
			{Index: 21, Name: "BirthPlace", DataType: "ST"},
			{Index: 22, Name: "MultipleBirthIndicator", DataType: "ID", Table: "HL70136"},
			{Index: 23, Name: "BirthOrder", DataType: "NM"},
			{Index: 24, Name: "Citizenship", DataType: "CE", Table: "HL70171"},
			{Index: 25, Name: "VeteransMilitaryStatus", DataType: "CE", Table: "HL70002"},
			{Index: 26, Name: "Nationality", DataType: "CE", Table: "HL70212"},
			{Index: 27, Name: "PatientDeathDateTime", DataType: "TS"},
			{Index: 28, Name: "PatientDeathIndicator", DataType: "ID", Table: "HL70136"},
		},
	},
	"PV1": {
		Name:          "PV1",
		IsRequired:    true,
		MinOccurrence: 1,
		MaxOccurrence: 1,
		Fields: []FieldDefinition{
			{Index: 1, Name: "SetID", DataType: "SI"},
			{Index: 2, Name: "PatientClass", DataType: "ID", IsRequired: true, Table: "HL70004"},
			{Index: 3, Name: "AssignedPatientLocation", DataType: "PL", Components: []FieldDefinition{
				{Index: 1, Name: "PointOfCare", DataType: "ST"},
				{Index: 2, Name: "Room", DataType: "ST"},
				{Index: 3, Name: "Bed", DataType: "ST"},
				{Index: 4, Name: "Facility", DataType: "HD"},
				{Index: 5, Name: "LocationStatus", DataType: "ST"},
				{Index: 6, Name: "PersonLocationType", DataType: "ST"},
			}},
			{Index: 4, Name: "AdmissionType", DataType: "ID", Table: "HL70007"},
			{Index: 5, Name: "PreadmitNumber", DataType: "CX"},
			{Index: 6, Name: "PriorPatientLocation", DataType: "PL"},
			{Index: 7, Name: "AttendingDoctor", DataType: "XCN", Components: []FieldDefinition{
				{Index: 1, Name: "IDNumber", DataType: "ST"},
				{Index: 2, Name: "FamilyName", DataType: "FN"},
				{Index: 3, Name: "GivenName", DataType: "ST"},
			}},
			{Index: 8, Name: "ReferringDoctor", DataType: "XCN"},
			{Index: 9, Name: "ConsultingDoctor", DataType: "XCN"},
			{Index: 10, Name: "HospitalService", DataType: "TS"},
			{Index: 11, Name: "TemporaryLocation", DataType: "PL"},
			{Index: 12, Name: "PreadmitTestIndicator", DataType: "ID", Table: "HL70087"},
			{Index: 13, Name: "ReAdmissionIndicator", DataType: "ID", Table: "HL70112"},
			{Index: 14, Name: "AdmitSource", DataType: "ID", Table: "HL70023"},
			{Index: 15, Name: "AmbulatoryStatus", DataType: "IS", Table: "HL70009"},
			{Index: 16, Name: "VIPIndicator", DataType: "ID", Table: "HL70036"},
			{Index: 17, Name: "AdmittingDoctor", DataType: "XCN"},
			{Index: 18, Name: "PatientType", DataType: "ID", Table: "HL70018"},
			{Index: 19, Name: "VisitNumber", DataType: "CX"},
			{Index: 20, Name: "FinancialClass", DataType: "FC"},
			{Index: 44, Name: "AdmitDateTime", DataType: "TS"},
			{Index: 45, Name: "DischargeDateTime", DataType: "TS"},
		},
	},
	"OBR": {
		Name:          "OBR",
		IsRequired:    false,
		MinOccurrence: 0,
		MaxOccurrence: -1,
		Fields: []FieldDefinition{
			{Index: 1, Name: "SetID", DataType: "SI"},
			{Index: 2, Name: "PlacerOrderNumber", DataType: "EI"},
			{Index: 3, Name: "FillerOrderNumber", DataType: "EI"},
			{Index: 4, Name: "UniversalServiceID", DataType: "CE", IsRequired: true},
			{Index: 5, Name: "Priority", DataType: "ID"},
			{Index: 6, Name: "RequestedDateTime", DataType: "TS"},
			{Index: 7, Name: "ObservationDateTime", DataType: "TS"},
			{Index: 8, Name: "ObservationEndDateTime", DataType: "TS"},
			{Index: 9, Name: "CollectionVolume", DataType: "CQ"},
			{Index: 10, Name: "CollectorIdentifier", DataType: "XCN"},
			{Index: 11, Name: "SpecimenActionCode", DataType: "ID", Table: "HL70065"},
			{Index: 12, Name: "DangerCode", DataType: "CE"},
			{Index: 13, Name: "RelevantClinicalInfo", DataType: "ST"},
			{Index: 14, Name: "SpecimenReceivedDateTime", DataType: "TS"},
			{Index: 15, Name: "SpecimenSource", DataType: "SPS"},
			{Index: 16, Name: "OrderingProvider", DataType: "XCN"},
			{Index: 17, Name: "OrderCallbackPhoneNumber", DataType: "XTN"},
			{Index: 18, Name: "PlacerField1", DataType: "ST"},
			{Index: 19, Name: "PlacerField2", DataType: "ST"},
			{Index: 20, Name: "FillerField1", DataType: "ST"},
			{Index: 21, Name: "FillerField2", DataType: "ST"},
			{Index: 22, Name: "ResultsRptStatusChangeDateTime", DataType: "TS"},
			{Index: 23, Name: "ChargeToPractice", DataType: "MOC"},
			{Index: 24, Name: "DiagnosticServSectionID", DataType: "ID", Table: "HL70074"},
			{Index: 25, Name: "ResultStatus", DataType: "ID", Table: "HL70123"},
			{Index: 26, Name: "ParentResult", DataType: "PRL"},
			{Index: 27, Name: "QuantityTiming", DataType: "TQ"},
			{Index: 28, Name: "ResultCopiesTo", DataType: "XCN"},
			{Index: 29, Name: "Parent", DataType: "EIP"},
			{Index: 30, Name: "TransportationMode", DataType: "ID", Table: "HL70124"},
			{Index: 31, Name: "ReasonForStudy", DataType: "CE"},
			{Index: 32, Name: "PrincipalResultInterpreter", DataType: "NDL"},
			{Index: 33, Name: "AssistantResultInterpreter", DataType: "NDL"},
			{Index: 34, Name: "Technician", DataType: "NDL"},
			{Index: 35, Name: "Transcriptionist", DataType: "NDL"},
			{Index: 36, Name: "ScheduledDateTime", DataType: "TS"},
			{Index: 37, Name: "NumberOfSampleContainers", DataType: "NM"},
			{Index: 38, Name: "TransportLogisticsOfSample", DataType: "CE"},
			{Index: 39, Name: "CollectorsComment", DataType: "CE"},
			{Index: 40, Name: "TransportArrangementResponsibility", DataType: "CE"},
			{Index: 41, Name: "TransportArranged", DataType: "ID", Table: "HL70224"},
			{Index: 42, Name: "EscortRequired", DataType: "ID", Table: "HL70225"},
			{Index: 43, Name: "PlannedPatientTransportComment", DataType: "CE"},
		},
	},
	"OBX": {
		Name:          "OBX",
		IsRequired:    false,
		MinOccurrence: 0,
		MaxOccurrence: -1,
		Fields: []FieldDefinition{
			{Index: 1, Name: "SetID", DataType: "SI"},
			{Index: 2, Name: "ValueType", DataType: "ID", Table: "HL70125"},
			{Index: 3, Name: "ObservationIdentifier", DataType: "CE", IsRequired: true, Components: []FieldDefinition{
				{Index: 1, Name: "Identifier", DataType: "ST"},
				{Index: 2, Name: "Text", DataType: "ST"},
				{Index: 3, Name: "CodingSystem", DataType: "ST"},
			}},
			{Index: 4, Name: "ObservationSubID", DataType: "ST"},
			{Index: 5, Name: "ObservationValue", DataType: "Varies", IsRequired: true},
			{Index: 6, Name: "Units", DataType: "CE"},
			{Index: 7, Name: "ReferenceRange", DataType: "ST"},
			{Index: 8, Name: "AbnormalFlags", DataType: "ID", Table: "HL70078"},
			{Index: 9, Name: "Probability", DataType: "NM"},
			{Index: 10, Name: "NatureOfAbnormalTest", DataType: "ID", Table: "HL70080"},
			{Index: 11, Name: "ResultStatus", DataType: "ID", Table: "HL70085"},
			{Index: 12, Name: "DateLastObsNormalValue", DataType: "TS"},
			{Index: 13, Name: "UserDefinedAccessChecks", DataType: "ST"},
			{Index: 14, Name: "DateTimeOfObservation", DataType: "TS"},
			{Index: 15, Name: "ProducersID", DataType: "CE"},
			{Index: 16, Name: "ResponsibleObserver", DataType: "XCN"},
			{Index: 17, Name: "ObservationMethod", DataType: "CE"},
			{Index: 18, Name: "EquipmentInstanceIdentifier", DataType: "EI"},
			{Index: 19, Name: "DateTimeOfAnalysis", DataType: "TS"},
			{Index: 20, Name: "Reserved", DataType: "ST"},
			{Index: 21, Name: "Reserved", DataType: "ST"},
			{Index: 22, Name: "Reserved", DataType: "ST"},
			{Index: 23, Name: "PerformingOrganizationName", DataType: "XON"},
			{Index: 24, Name: "PerformingOrganizationAddress", DataType: "XAD"},
			{Index: 25, Name: "PerformingOrganizationMedicalDirector", DataType: "XCN"},
		},
	},
	"NK1": {
		Name:          "NK1",
		IsRequired:    false,
		MinOccurrence: 0,
		MaxOccurrence: -1,
		Fields: []FieldDefinition{
			{Index: 1, Name: "SetID", DataType: "SI"},
			{Index: 2, Name: "Name", DataType: "XPN"},
			{Index: 3, Name: "Relationship", DataType: "CE", Table: "HL70063"},
			{Index: 4, Name: "Address", DataType: "XAD"},
			{Index: 5, Name: "PhoneNumber", DataType: "XTN"},
			{Index: 6, Name: "BusinessPhoneNumber", DataType: "XTN"},
			{Index: 7, Name: "ContactRole", DataType: "CE", Table: "HL70131"},
			{Index: 8, Name: "StartDate", DataType: "TS"},
			{Index: 9, Name: "EndDate", DataType: "TS"},
			{Index: 10, Name: "NextOfKinAddresses", DataType: "XAD"},
			{Index: 11, Name: "NextOfKinPhoneNumbers", DataType: "XTN"},
		},
	},
	"DG1": {
		Name:          "DG1",
		IsRequired:    false,
		MinOccurrence: 0,
		MaxOccurrence: -1,
		Fields: []FieldDefinition{
			{Index: 1, Name: "SetID", DataType: "SI"},
			{Index: 2, Name: "DiagnosisCodingMethod", DataType: "ID", Table: "HL70053"},
			{Index: 3, Name: "DiagnosisCode", DataType: "CE", IsRequired: true, Table: "HL70051"},
			{Index: 4, Name: "DiagnosisDescription", DataType: "ST"},
			{Index: 5, Name: "DiagnosisDateTime", DataType: "TS"},
			{Index: 6, Name: "DiagnosisType", DataType: "ID", IsRequired: true, Table: "HL70052"},
		},
	},
	"AL1": {
		Name:          "AL1",
		IsRequired:    false,
		MinOccurrence: 0,
		MaxOccurrence: -1,
		Fields: []FieldDefinition{
			{Index: 1, Name: "SetID", DataType: "SI"},
			{Index: 2, Name: "AllergenTypeCode", DataType: "CE", Table: "HL70027"},
			{Index: 3, Name: "AllergenCode", DataType: "CE", IsRequired: true},
			{Index: 4, Name: "AllergySeverityCode", DataType: "CE", Table: "HL70028"},
			{Index: 5, Name: "AllergyReactionCode", DataType: "ST"},
			{Index: 6, Name: "IdentificationDate", DataType: "TS"},
		},
	},
	"GT1": {
		Name:          "GT1",
		IsRequired:    false,
		MinOccurrence: 0,
		MaxOccurrence: -1,
		Fields: []FieldDefinition{
			{Index: 1, Name: "SetID", DataType: "SI"},
			{Index: 2, Name: "GuarantorNumber", DataType: "CX"},
			{Index: 3, Name: "GuarantorName", DataType: "XPN", IsRequired: true},
			{Index: 4, Name: "GuarantorSpouseName", DataType: "XPN"},
			{Index: 5, Name: "GuarantorAddress", DataType: "XAD"},
			{Index: 6, Name: "GuarantorHomePhone", DataType: "XTN"},
			{Index: 7, Name: "GuarantorBusinessPhone", DataType: "XTN"},
			{Index: 8, Name: "GuarantorDateOfBirth", DataType: "TS"},
			{Index: 9, Name: "GuarantorSex", DataType: "ID", Table: "HL70001"},
			{Index: 10, Name: "GuarantorType", DataType: "IS", Table: "HL70068"},
			{Index: 11, Name: "GuarantorRelationship", DataType: "CE", Table: "HL70063"},
			{Index: 12, Name: "GuarantorSSN", DataType: "ST"},
			{Index: 13, Name: "GuarantorDateBegin", DataType: "TS"},
			{Index: 14, Name: "GuarantorDateEnd", DataType: "TS"},
			{Index: 15, Name: "GuarantorPriority", DataType: "NM"},
		},
	},
	"IN1": {
		Name:          "IN1",
		IsRequired:    false,
		MinOccurrence: 0,
		MaxOccurrence: -1,
		Fields: []FieldDefinition{
			{Index: 1, Name: "SetID", DataType: "SI"},
			{Index: 2, Name: "InsurancePlanID", DataType: "CE", IsRequired: true},
			{Index: 3, Name: "InsuranceCompanyID", DataType: "CX"},
			{Index: 4, Name: "InsuranceCompanyName", DataType: "XON"},
			{Index: 5, Name: "InsuranceCompanyAddress", DataType: "XAD"},
			{Index: 6, Name: "InsuranceCoContactPerson", DataType: "XPN"},
			{Index: 7, Name: "InsuranceCoPhoneNumber", DataType: "XTN"},
			{Index: 8, Name: "GroupNumber", DataType: "ST"},
			{Index: 9, Name: "GroupName", DataType: "XON"},
			{Index: 10, Name: "PlanEffectiveDate", DataType: "TS"},
			{Index: 11, Name: "PlanExpirationDate", DataType: "TS"},
			{Index: 12, Name: "AuthorizationInformation", DataType: "AUI"},
			{Index: 13, Name: "PlanType", DataType: "IS", Table: "HL70086"},
			{Index: 14, Name: "NameOfInsured", DataType: "XPN"},
			{Index: 15, Name: "InsuredRelationship", DataType: "CE", Table: "HL70063"},
			{Index: 16, Name: "InsuredDateOfBirth", DataType: "TS"},
			{Index: 17, Name: "InsuredAddress", DataType: "XAD"},
			{Index: 18, Name: "AssignmentOfBenefits", DataType: "ID", Table: "HL70135"},
			{Index: 19, Name: "CoordinationOfBenefits", DataType: "ID", Table: "HL70136"},
			{Index: 20, Name: "CoordBenPriority", DataType: "ST"},
			{Index: 21, Name: "NoticeOfAdmissionCode", DataType: "ID", Table: "HL70136"},
			{Index: 22, Name: "NoticeOfAdmissionDate", DataType: "TS"},
			{Index: 23, Name: "ReportOfEligibilityCode", DataType: "ID", Table: "HL70136"},
			{Index: 24, Name: "ReportOfEligibilityDate", DataType: "TS"},
			{Index: 25, Name: "ReleaseInformationCode", DataType: "IS", Table: "HL70093"},
			{Index: 26, Name: "DelayDays", DataType: "ST"},
			{Index: 27, Name: "PlanID", DataType: "CE"},
			{Index: 28, Name: "InvalidReason", DataType: "ID", Table: "HL70143"},
			{Index: 29, Name: "EffectiveDate", DataType: "TS"},
			{Index: 30, Name: "ExpirationDate", DataType: "TS"},
		},
	},
	"NTE": {
		Name:          "NTE",
		IsRequired:    false,
		MinOccurrence: 0,
		MaxOccurrence: -1,
		Fields: []FieldDefinition{
			{Index: 1, Name: "SetID", DataType: "SI"},
			{Index: 2, Name: "SourceOfComment", DataType: "ID", Table: "HL70105"},
			{Index: 3, Name: "Comment", DataType: "FT", MaxLength: 65536},
			{Index: 4, Name: "CommentType", DataType: "CE", Table: "HL70204"},
		},
	},
	"ORC": {
		Name:          "ORC",
		IsRequired:    false,
		MinOccurrence: 0,
		MaxOccurrence: 1,
		Fields: []FieldDefinition{
			{Index: 1, Name: "OrderControl", DataType: "ID", IsRequired: true, Table: "HL70119"},
			{Index: 2, Name: "PlacerOrderNumber", DataType: "EI"},
			{Index: 3, Name: "FillerOrderNumber", DataType: "EI"},
			{Index: 4, Name: "PlacerGroupNumber", DataType: "EI"},
			{Index: 5, Name: "OrderStatus", DataType: "ID", Table: "HL70038"},
			{Index: 6, Name: "ResponseFlag", DataType: "ID", Table: "HL70121"},
			{Index: 7, Name: "QuantityTiming", DataType: "TQ"},
			{Index: 8, Name: "Parent", DataType: "EIP"},
			{Index: 9, Name: "DateTimeOfTransaction", DataType: "TS"},
			{Index: 10, Name: "EnteredBy", DataType: "XCN"},
			{Index: 11, Name: "VerifiedBy", DataType: "XCN"},
			{Index: 12, Name: "OrderingProvider", DataType: "XCN"},
			{Index: 13, Name: "EntererLocation", DataType: "PL"},
			{Index: 14, Name: "CallBackPhoneNumber", DataType: "XTN"},
			{Index: 15, Name: "OrderEffectiveDateTime", DataType: "TS"},
			{Index: 16, Name: "OrderControlCodeReason", DataType: "CE"},
			{Index: 17, Name: "EnteringOrganization", DataType: "CE"},
			{Index: 18, Name: "EnteringDevice", DataType: "CE"},
			{Index: 19, Name: "ActionBy", DataType: "XCN"},
		},
	},
}

type MessageStructure struct {
	MessageType string
	Version     string
	Segments    []SegmentRequirement
}

type SegmentRequirement struct {
	Name          string
	SegmentType   string
	Position      int
	Group         string
	IsRequired    bool
	MaxOccurrence int
}

var MessageStructures = map[string]MessageStructure{
	"ADT_A01": {
		MessageType: "ADT^A01",
		Segments: []SegmentRequirement{
			{Name: "MSH", Position: 1, IsRequired: true},
			{Name: "SFT", Position: 2, IsRequired: false},
			{Name: "EVN", Position: 3, IsRequired: true},
			{Name: "PID", Position: 4, IsRequired: true},
			{Name: "PD1", Position: 5, IsRequired: false},
			{Name: "ROL", Position: 6, IsRequired: false},
			{Name: "NK1", Position: 7, IsRequired: false, MaxOccurrence: -1},
			{Name: "PV1", Position: 8, IsRequired: true},
			{Name: "PV2", Position: 9, IsRequired: false},
			{Name: "ROL", Position: 10, IsRequired: false},
			{Name: "DB1", Position: 11, IsRequired: false},
			{Name: "OBX", Position: 12, IsRequired: false, MaxOccurrence: -1},
			{Name: "AL1", Position: 13, IsRequired: false, MaxOccurrence: -1},
			{Name: "DG1", Position: 14, IsRequired: false, MaxOccurrence: -1},
			{Name: "DRG", Position: 15, IsRequired: false},
			{Name: "PR1", Position: 16, IsRequired: false, MaxOccurrence: -1},
			{Name: "GT1", Position: 17, IsRequired: false, MaxOccurrence: -1},
			{Name: "IN1", Position: 18, IsRequired: false, MaxOccurrence: -1},
			{Name: "IN2", Position: 19, IsRequired: false},
			{Name: "IN3", Position: 20, IsRequired: false, MaxOccurrence: -1},
			{Name: "ACC", Position: 21, IsRequired: false},
			{Name: "UB1", Position: 22, IsRequired: false},
			{Name: "UB2", Position: 23, IsRequired: false},
		},
	},
	"ADT_A04": {
		MessageType: "ADT^A04",
		Segments: []SegmentRequirement{
			{Name: "MSH", Position: 1, IsRequired: true},
			{Name: "SFT", Position: 2, IsRequired: false},
			{Name: "EVN", Position: 3, IsRequired: true},
			{Name: "PID", Position: 4, IsRequired: true},
			{Name: "PD1", Position: 5, IsRequired: false},
			{Name: "ROL", Position: 6, IsRequired: false},
			{Name: "NK1", Position: 7, IsRequired: false, MaxOccurrence: -1},
			{Name: "PV1", Position: 8, IsRequired: true},
			{Name: "PV2", Position: 9, IsRequired: false},
			{Name: "ROL", Position: 10, IsRequired: false},
			{Name: "DB1", Position: 11, IsRequired: false},
			{Name: "OBX", Position: 12, IsRequired: false, MaxOccurrence: -1},
			{Name: "AL1", Position: 13, IsRequired: false, MaxOccurrence: -1},
			{Name: "DG1", Position: 14, IsRequired: false, MaxOccurrence: -1},
			{Name: "DRG", Position: 15, IsRequired: false},
		},
	},
	"ORU_R01": {
		MessageType: "ORU^R01",
		Segments: []SegmentRequirement{
			{Name: "MSH", Position: 1, IsRequired: true},
			{Name: "SFT", Position: 2, IsRequired: false},
			{Name: "NTE", Position: 3, IsRequired: false, MaxOccurrence: -1},
			{Name: "PATIENT_RESULT", Position: 4, IsRequired: true, MaxOccurrence: -1},
		},
	},
}

type HL7Table struct {
	TableID   string
	TableName string
	Values    map[string]string
}

var HL7Tables = map[string]HL7Table{
	"HL70001": {
		TableID:   "HL70001",
		TableName: "Administrative Sex",
		Values: map[string]string{
			"A": "Ambiguous",
			"F": "Female",
			"M": "Male",
			"N": "Not applicable",
			"O": "Other",
			"U": "Unknown",
		},
	},
	"HL70002": {
		TableID:   "HL70002",
		TableName: "Marital Status",
		Values: map[string]string{
			"A": "Separated",
			"D": "Divorced",
			"I": "Interlocutory",
			"M": "Married",
			"P": "Polygamous",
			"S": "Single",
			"T": "Domestic partner",
			"U": "Unknown",
			"W": "Widowed",
		},
	},
	"HL70004": {
		TableID:   "HL70004",
		TableName: "Patient Class",
		Values: map[string]string{
			"B": "Obstetrics",
			"C": "Commercial Account",
			"E": "Emergency",
			"I": "Invalid",
			"N": "Not Applicable",
			"O": "Outpatient",
			"P": "Preadmit",
			"R": "Recurring patient",
			"U": "Unknown",
		},
	},
	"HL70006": {
		TableID:   "HL70006",
		TableName: "Religion",
		Values: map[string]string{
			"AG":  "Anglican",
			"BAP": "Baptist",
			"BTH": "Buddhist",
			"CAT": "Roman Catholic",
			"CHM": "Christian",
			"CON": "Confucian",
			"DOC": "Doctor of Christianity",
			"E":   "Episcopal",
			"EMC": "Eastern Orthodox",
			"ETH": "Ethiopian Orthodox",
			"F":   "Christian (none of the above)",
			"FR":  "French Reformed",
			"FRE": "Friends",
			"G":   "Greek Orthodox",
			"H":   "Hindu",
			"JE":  "Jewish",
			"L":   "Lutheran",
			"MEN": "Mennonite",
			"MET": "Methodist",
			"MOS": "Mormon",
			"MU":  "Muslim",
			"N":   "None",
			"NON": "Nonreligious",
			"ORT": "Orthodox",
			"P":   "Presbyterian",
			"PA":  "Pagan",
			"PRC": "Other Christian",
			"PRO": "Protestant",
			"RE":  "Reformed",
			"REC": "Reformed Church",
			"RF":  "Reformed",
			"S":   "Spiritist",
			"SAL": "Salvation Army",
			"SAM": "Seventh Adventist",
			"SD":  "Seventh Day Adventist",
			"SE":  "Sikh",
			"SH":  "Shintoist",
			"STA": "Storefront",
			"TY":  "Taoist",
			"U":   "Unknown",
			"UNC": "Unchurched",
			"UNF": "Unitarian",
			"VAR": "Unknown",
			"W":   "Wesleyan",
			"WEL": "Welsh",
			"Z":   "Zoroastrian",
		},
	},
	"HL70038": {
		TableID:   "HL70038",
		TableName: "Order Status",
		Values: map[string]string{
			"A":  "Some, but not all, results available",
			"CA": "Order was canceled",
			"CM": "Order is completed",
			"DC": "Order was discontinued",
			"ER": "Error, order cannot be processed",
			"HD": "Order is on hold",
			"IP": "In process",
			"RP": "Order has been replaced",
			"SC": "Order is scheduled",
		},
	},
	"HL70123": {
		TableID:   "HL70123",
		TableName: "Observation Result Status",
		Values: map[string]string{
			"C": "Record coming over is a correction and thus replaces a final result",
			"D": "Deletes the OBX record",
			"F": "Final result; can only be changed with a corrected result",
			"I": "Specimen in instrument; processing",
			"N": "Not asked; used to indicate no result was considered for this order",
			"O": "Order detail description only (no result)",
			"P": "Preliminary: a result that has been verified",
			"R": "Requested: result is not yet available",
			"S": "Partial: a preliminary result that is currently being processed",
			"U": "Results status not available",
			"W": "Post original as wrong",
			"X": "No results available; order was canceled",
		},
	},
	"HL70125": {
		TableID:   "HL70125",
		TableName: "Value Type",
		Values: map[string]string{
			"AD":   "Address",
			"CE":   "Coded Entry",
			"CF":   "Coded Element With Formatted Values",
			"CK":   "Composite ID With Check Digit",
			"CN":   "Composite ID And Name",
			"CNE":  "Coded With No Exceptions",
			"CNP":  "Coded With No Person Name",
			"CP":   "Composite Price",
			"CPN":  "Composite Person Name",
			"CS":   "Coded Simple Value",
			"CT":   "Coded Token",
			"CX":   "Extended Composite ID With Check Digit",
			"DD":   "Structured Numeric",
			"DLN":  "Driver's License Number",
			"DLT":  "Delta",
			"DR":   "Date/Time Range",
			"DT":   "Date",
			"ED":   "Encapsulated Data",
			"EI":   "Entity Identifier",
			"ELM":  "Encapsulated Data",
			"EN":   "Entity Name",
			"ERN":  "Entity Name",
			"FC":   "Financial Class",
			"FN":   "Family Name",
			"FT":   "Formatted Text",
			"FX":   "Fraction",
			"GTS":  "General Timing Specification",
			"HD":   "Hierarchic Designation",
			"HED":  "HL7 Encapsulated Data",
			"HI":   "Hierarchic Interface",
			"ICD":  "Insurance",
			"ID":   "Coded Value For HL7 Defined Tables",
			"IE":   "Identifier Expressive",
			"IIM":  "Inverse Identifier With Mandatory",
			"IL":   "Instance Identifier",
			"IS":   "Coded Value For User-Defined Tables",
			"JCC":  "Job Code Class",
			"LA":   "Location Address (using Street Address Line)",
			"LA1":  "Location Address",
			"LAH":  "Location Address",
			"LD":   "Location With Address",
			"LI":   "Long Integers",
			"LN":   "License Number",
			"LP":   "Location With Point Of Coordinates",
			"MC":   "Money",
			"MD":   "Multiple Types",
			"MF":   "Money And Percentage",
			"MOC":  "Charge Rate And Time",
			"MO":   "Money",
			"MOP":  "Money Or Percentage",
			"MSA":  "Message Acknowledgment",
			"MSG":  "Message Type",
			"NA":   "Numeric Array",
			"ND":   "Numeric With Decimal",
			"NI":   "National Identifier",
			"NM":   "Numeric",
			"NM1":  "Numeric",
			"NPI":  "National Provider Identifier",
			"NT":   "Numeric Token",
			"NU":   "Numeric",
			"OBR":  "Observation Request",
			"OD":   "Other Numeric",
			"OS":   "Other Scalar",
			"OUI":  "Operator Universal Identifier",
			"OW":   "Order Number And Date",
			"PN":   "Person Name",
			"PPN":  "Performing Person Name",
			"PR":   "Priority",
			"PRA":  "Practitioner ID",
			"PT":   "Processing Type",
			"PTA":  "Payment Type",
			"QIP":  "Query Input Parameter List",
			"QSC":  "Query Selection Criteria",
			"RCD":  "Row Column Definition",
			"RC":   "Regional Center Code",
			"RFR":  "Reference Range",
			"RI":   "Interval",
			"RMC":  "Room Coverage",
			"RP":   "Reference Pointer",
			"RPL":  "Reference Pointer To Location",
			"RP1":  "Referral Information",
			"RP2":  "Referral Information",
			"RP3":  "Referral Information",
			"RP5":  "Referral Information",
			"RPO":  "Referral And Order",
			"RNA":  "Rest Numeric Array",
			"RND":  "Rest Numeric",
			"SNM":  "String Of Name",
			"SN":   "Structured Numeric",
			"SN1":  "Structured Numeric",
			"SPD":  "Specialty",
			"SPS":  "Specimen Source",
			"ST":   "String Data",
			"ST1":  "String Data",
			"TM":   "Time",
			"TN":   "Telephone Number",
			"TQ":   "Timing Quantity",
			"TS":   "Time Stamp",
			"TX":   "Text Data",
			"UC":   "Urgency Classification",
			"UD":   "Undefined",
			"UI":   "Unique Identifier",
			"UID":  "Unique Identifier",
			"UPIN": "UPIN",
			"URI":  "Universal Resource Identifier",
			"UV":   "Undefined Variable",
			"V24":  "Value",
			"VAR":  "Variance",
			"VH":   "Verifying Healthcare Interest",
			"VID":  "Version Identifier",
			"VR":   "Value Range",
			"VT":   "Virtual Table",
			"WCD":  "Wrong Check Digit",
			"XAD":  "Extended Address",
			"XCN":  "Extended Composite ID Number And Name For Persons",
			"XON":  "Extended Composite Name And Identification Number For Organizations",
			"XPN":  "Extended Person Name",
			"XT":   "Text",
			"XTN":  "Extended Telecommunications Number",
		},
	},
	"HL70305": {
		TableID:   "HL70305",
		TableName: "Diagnosis Coding Method",
		Values: map[string]string{
			"01": "ICD-9",
			"02": "ICD-10",
			"99": "Other",
		},
	},
	"HL70394": {
		TableID:   "HL70394",
		TableName: "Message Error Condition Codes",
		Values: map[string]string{
			"0":   "Message accepted",
			"100": "Segment sequence error",
			"101": "Required field missing",
			"102": "Data type error",
			"103": "Table value not found",
			"104": "Unsupported message type",
			"105": "Unsupported event code",
			"106": "Unsupported processing ID",
			"107": "Unsupported version ID",
			"108": "Unknown key identifier",
			"109": "Duplicate key identifier",
			"110": "Application record error",
			"111": "Field content error",
			"112": "Segmentation error",
			"200": "Unsupported character set",
			"207": "Application internal error",
		},
	},
}

type (
	PersonName struct {
		FamilyName             string `hl7:"1"`
		GivenName              string `hl7:"2"`
		SecondAndFurtherNames  string `hl7:"3"`
		Suffix                 string `hl7:"4"`
		Prefix                 string `hl7:"5"`
		Degree                 string `hl7:"6"`
		NameTypeCode           string `hl7:"7"`
		NameRepresentationCode string `hl7:"8"`
	}

	CodedElement struct {
		Identifier      string `hl7:"1"`
		Text            string `hl7:"2"`
		CodingSystem    string `hl7:"3"`
		AltIdentifier   string `hl7:"4"`
		AltText         string `hl7:"5"`
		AltCodingSystem string `hl7:"6"`
	}

	ExtendedAddress struct {
		StreetAddress              string `hl7:"1"`
		OtherDesignation           string `hl7:"2"`
		City                       string `hl7:"3"`
		StateOrProvince            string `hl7:"4"`
		ZipOrPostalCode            string `hl7:"5"`
		Country                    string `hl7:"6"`
		AddressType                string `hl7:"7"`
		OtherGeographicDesignation string `hl7:"8"`
		CountyCode                 string `hl7:"9"`
		AddressRepresentationCode  string `hl7:"10"`
	}

	ExtendedPhoneNumber struct {
		TelephoneNumber          string `hl7:"1"`
		TelecommunicationUseCode string `hl7:"2"`
		EquipmentType            string `hl7:"3"`
		EmailAddress             string `hl7:"4"`
		CountryCode              string `hl7:"5"`
		AreaCityCode             string `hl7:"6"`
		LocalNumber              string `hl7:"7"`
		Extension                string `hl7:"8"`
		Text                     string `hl7:"9"`
	}

	ExtendedCompositeID struct {
		IDNumber           string `hl7:"1"`
		CheckDigit         string `hl7:"2"`
		CheckDigitScheme   string `hl7:"3"`
		AssigningAuthority string `hl7:"4"`
		IdentifierTypeCode string `hl7:"5"`
		AssigningFacility  string `hl7:"6"`
		EffectiveDate      string `hl7:"7"`
		ExpirationDate     string `hl7:"8"`
	}

	Component struct {
		Time      time.Time
		Timezone  string
		Precision string
	}
)

func LookupTable(tableID, code string) (string, bool) {
	if table, ok := HL7Tables[tableID]; ok {
		if desc, ok := table.Values[code]; ok {
			return desc, true
		}
	}
	return "", false
}

func LookupSegmentDefinition(segmentName string) (SegmentDefinition, bool) {
	if def, ok := SegmentDefinitions[segmentName]; ok {
		return def, true
	}
	return SegmentDefinition{}, false
}

func LookupMessageStructure(messageType string) (MessageStructure, bool) {
	if ms, ok := MessageStructures[messageType]; ok {
		return ms, true
	}
	for _, ms := range MessageStructures {
		if ms.MessageType == messageType {
			return ms, true
		}
	}
	return MessageStructure{}, false
}

func GetFieldDefinition(segmentName string, fieldIndex int) (FieldDefinition, bool) {
	if def, ok := SegmentDefinitions[segmentName]; ok {
		for _, field := range def.Fields {
			if field.Index == fieldIndex {
				return field, true
			}
		}
	}
	return FieldDefinition{}, false
}
