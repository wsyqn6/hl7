package hl7

import "time"

type PIDHelper struct {
	seg Segment
}

func (m *Message) PID() PIDHelper {
	if seg, ok := m.Segment("PID"); ok {
		return PIDHelper{seg: seg}
	}
	return PIDHelper{}
}

func (h PIDHelper) Exists() bool {
	return h.seg.Name() != ""
}

func (h PIDHelper) PatientID() string {
	return h.seg.Field(3)
}

func (h PIDHelper) PatientIdentifierList() []string {
	return SplitField(h.seg.Field(3), '^')
}

func (h PIDHelper) AlternatePatientID() string {
	return h.seg.Field(4)
}

func (h PIDHelper) PatientName() string {
	return h.seg.Field(5)
}

func (h PIDHelper) LastName() string {
	return ParseComponent(h.seg.Field(5), 1)
}

func (h PIDHelper) FirstName() string {
	return ParseComponent(h.seg.Field(5), 2)
}

func (h PIDHelper) MiddleName() string {
	return ParseComponent(h.seg.Field(5), 3)
}

func (h PIDHelper) Suffix() string {
	return ParseComponent(h.seg.Field(5), 4)
}

func (h PIDHelper) Prefix() string {
	return ParseComponent(h.seg.Field(5), 5)
}

func (h PIDHelper) Degree() string {
	return ParseComponent(h.seg.Field(5), 6)
}

func (h PIDHelper) MothersMaidenName() string {
	return h.seg.Field(6)
}

func (h PIDHelper) DateOfBirth() string {
	return h.seg.Field(7)
}

func (h PIDHelper) DateOfBirthTime() (time.Time, error) {
	return time.Parse("20060102150405", h.seg.Field(7))
}

func (h PIDHelper) Sex() string {
	return h.seg.Field(8)
}

func (h PIDHelper) Gender() string {
	return h.seg.Field(8)
}

func (h PIDHelper) PatientAlias() string {
	return h.seg.Field(9)
}

func (h PIDHelper) Race() string {
	return h.seg.Field(10)
}

func (h PIDHelper) PatientAddress() string {
	return h.seg.Field(11)
}

func (h PIDHelper) StreetAddress() string {
	return ParseComponent(h.seg.Field(11), 1)
}

func (h PIDHelper) StreetName() string {
	return ParseComponent(h.seg.Field(11), 2)
}

func (h PIDHelper) City() string {
	return ParseComponent(h.seg.Field(11), 3)
}

func (h PIDHelper) State() string {
	return ParseComponent(h.seg.Field(11), 4)
}

func (h PIDHelper) PostalCode() string {
	return ParseComponent(h.seg.Field(11), 5)
}

func (h PIDHelper) Country() string {
	return ParseComponent(h.seg.Field(11), 6)
}

func (h PIDHelper) PhoneHome() string {
	return h.seg.Field(13)
}

func (h PIDHelper) PhoneBusiness() string {
	return h.seg.Field(14)
}

func (h PIDHelper) PrimaryLanguage() string {
	return h.seg.Field(15)
}

func (h PIDHelper) MaritalStatus() string {
	return h.seg.Field(16)
}

func (h PIDHelper) Religion() string {
	return h.seg.Field(17)
}

func (h PIDHelper) SSN() string {
	return h.seg.Field(19)
}

func (h PIDHelper) DriversLicense() string {
	return h.seg.Field(20)
}

func (h PIDHelper) EthnicGroup() string {
	return h.seg.Field(22)
}

func (h PIDHelper) BirthPlace() string {
	return h.seg.Field(23)
}

func (h PIDHelper) MultipleBirthIndicator() string {
	return h.seg.Field(24)
}

func (h PIDHelper) BirthOrder() string {
	return h.seg.Field(25)
}

func (h PIDHelper) Citizenship() []string {
	return SplitField(h.seg.Field(26), '^')
}

func (h PIDHelper) VeteransMilitaryStatus() string {
	return h.seg.Field(27)
}

type MSHHelper struct {
	seg Segment
}

func (m *Message) MSH() MSHHelper {
	if seg, ok := m.Segment("MSH"); ok {
		return MSHHelper{seg: seg}
	}
	return MSHHelper{}
}

func (h MSHHelper) Exists() bool {
	return h.seg.Name() != ""
}

func (h MSHHelper) FieldSeparator() string {
	return h.seg.Field(1)
}

func (h MSHHelper) EncodingCharacters() string {
	return h.seg.Field(2)
}

func (h MSHHelper) SendingApplication() string {
	return h.seg.Field(3)
}

func (h MSHHelper) SendingFacility() string {
	return h.seg.Field(4)
}

func (h MSHHelper) ReceivingApplication() string {
	return h.seg.Field(5)
}

func (h MSHHelper) ReceivingFacility() string {
	return h.seg.Field(6)
}

func (h MSHHelper) DateTimeOfMessage() string {
	return h.seg.Field(7)
}

func (h MSHHelper) DateTime() time.Time {
	t, _ := time.Parse("20060102150405", h.seg.Field(7))
	return t
}

func (h MSHHelper) Security() string {
	return h.seg.Field(8)
}

func (h MSHHelper) MessageType() string {
	return h.seg.Field(9)
}

func (h MSHHelper) MessageTypeCode() string {
	return ParseComponent(h.seg.Field(9), 1)
}

func (h MSHHelper) MessageTypeTrigger() string {
	return ParseComponent(h.seg.Field(9), 2)
}

func (h MSHHelper) MessageTypeStructure() string {
	return ParseComponent(h.seg.Field(9), 3)
}

func (h MSHHelper) MessageControlID() string {
	return h.seg.Field(10)
}

func (h MSHHelper) ProcessingID() string {
	return h.seg.Field(11)
}

func (h MSHHelper) VersionID() string {
	return h.seg.Field(12)
}

func (h MSHHelper) SequenceNumber() string {
	return h.seg.Field(13)
}

func (h MSHHelper) ContinuationPointer() string {
	return h.seg.Field(14)
}

func (h MSHHelper) AcceptAcknowledgmentType() string {
	return h.seg.Field(15)
}

func (h MSHHelper) ApplicationAcknowledgmentType() string {
	return h.seg.Field(16)
}

func (h MSHHelper) CountryCode() string {
	return h.seg.Field(17)
}

type PV1Helper struct {
	seg Segment
}

func (m *Message) PV1() PV1Helper {
	if seg, ok := m.Segment("PV1"); ok {
		return PV1Helper{seg: seg}
	}
	return PV1Helper{}
}

func (h PV1Helper) Exists() bool {
	return h.seg.Name() != ""
}

func (h PV1Helper) SetID() string {
	return h.seg.Field(1)
}

func (h PV1Helper) PatientClass() string {
	return h.seg.Field(2)
}

func (h PV1Helper) AssignedPatientLocation() string {
	return h.seg.Field(3)
}

func (h PV1Helper) LocationPointOfCare() string {
	return ParseComponent(h.seg.Field(3), 1)
}

func (h PV1Helper) LocationRoom() string {
	return ParseComponent(h.seg.Field(3), 2)
}

func (h PV1Helper) LocationBed() string {
	return ParseComponent(h.seg.Field(3), 3)
}

func (h PV1Helper) LocationFacility() string {
	return ParseComponent(h.seg.Field(3), 4)
}

func (h PV1Helper) AdmissionType() string {
	return h.seg.Field(4)
}

func (h PV1Helper) AttendingDoctor() string {
	return h.seg.Field(7)
}

func (h PV1Helper) AttendingDoctorID() string {
	return ParseComponent(h.seg.Field(7), 1)
}

func (h PV1Helper) AttendingDoctorName() string {
	return ParseComponent(h.seg.Field(7), 2)
}

func (h PV1Helper) ReferringDoctor() string {
	return h.seg.Field(8)
}

func (h PV1Helper) ConsultingDoctor() []string {
	var doctors []string
	for i := 9; i <= 12; i++ {
		if d := h.seg.Field(i); d != "" {
			doctors = append(doctors, d)
		}
	}
	return doctors
}

func (h PV1Helper) HospitalService() string {
	return h.seg.Field(10)
}

func (h PV1Helper) AdmitSource() string {
	return h.seg.Field(13)
}

func (h PV1Helper) VIPIndicator() string {
	return h.seg.Field(16)
}

func (h PV1Helper) AdmissionDate() string {
	return h.seg.Field(44)
}

func (h PV1Helper) DischargeDate() string {
	return h.seg.Field(45)
}

type OBRHelper struct {
	seg Segment
}

func (m *Message) OBR() OBRHelper {
	if seg, ok := m.Segment("OBR"); ok {
		return OBRHelper{seg: seg}
	}
	return OBRHelper{}
}

func (h OBRHelper) Exists() bool {
	return h.seg.Name() != ""
}

func (h OBRHelper) SetID() string {
	return h.seg.Field(1)
}

func (h OBRHelper) PlacerOrderNumber() string {
	return h.seg.Field(2)
}

func (h OBRHelper) FillerOrderNumber() string {
	return h.seg.Field(3)
}

func (h OBRHelper) UniversalServiceID() string {
	return h.seg.Field(4)
}

func (h OBRHelper) ServiceIdentifier() string {
	return ParseComponent(h.seg.Field(4), 1)
}

func (h OBRHelper) ServiceText() string {
	return ParseComponent(h.seg.Field(4), 2)
}

func (h OBRHelper) Priority() string {
	return h.seg.Field(5)
}

func (h OBRHelper) RequestedDateTime() string {
	return h.seg.Field(6)
}

func (h OBRHelper) ObservationDateTime() string {
	return h.seg.Field(7)
}

func (h OBRHelper) ObservationEndDateTime() string {
	return h.seg.Field(8)
}

func (h OBRHelper) CollectionVolume() string {
	return h.seg.Field(9)
}

func (h OBRHelper) CollectorIdentifier() string {
	return h.seg.Field(10)
}

func (h OBRHelper) SpecimenActionCode() string {
	return h.seg.Field(11)
}

func (h OBRHelper) DangerCode() string {
	return h.seg.Field(12)
}

func (h OBRHelper) RelevantClinicalInfo() string {
	return h.seg.Field(13)
}

func (h OBRHelper) SpecimenReceivedDateTime() string {
	return h.seg.Field(14)
}

func (h OBRHelper) SpecimenSource() string {
	return h.seg.Field(15)
}

func (h OBRHelper) OrderingProvider() string {
	return h.seg.Field(16)
}

func (h OBRHelper) OrderCallbackPhoneNumber() string {
	return h.seg.Field(17)
}

func (h OBRHelper) FillerStatusCode() string {
	return h.seg.Field(25)
}

func (h OBRHelper) ResultsRptStatusChangeDateTime() string {
	return h.seg.Field(22)
}

func (h OBRHelper) TransportationMode() string {
	return h.seg.Field(23)
}

func (h OBRHelper) ReasonForStudy() string {
	return h.seg.Field(31)
}

type OBXHelper struct {
	seg Segment
}

func (m *Message) OBX() OBXHelper {
	if seg, ok := m.Segment("OBX"); ok {
		return OBXHelper{seg: seg}
	}
	return OBXHelper{}
}

func (h OBXHelper) Exists() bool {
	return h.seg.Name() != ""
}

func (h OBXHelper) SetID() string {
	return h.seg.Field(1)
}

func (h OBXHelper) ValueType() string {
	return h.seg.Field(2)
}

func (h OBXHelper) ObservationIdentifier() string {
	return h.seg.Field(3)
}

func (h OBXHelper) ObservationIdentifierCode() string {
	return ParseComponent(h.seg.Field(3), 1)
}

func (h OBXHelper) ObservationIdentifierText() string {
	return ParseComponent(h.seg.Field(3), 2)
}

func (h OBXHelper) ObservationSubID() string {
	return h.seg.Field(4)
}

func (h OBXHelper) ObservationValue() string {
	return h.seg.Field(5)
}

func (h OBXHelper) Units() string {
	return h.seg.Field(6)
}

func (h OBXHelper) ReferenceRange() string {
	return h.seg.Field(7)
}

func (h OBXHelper) AbnormalFlags() []string {
	return SplitField(h.seg.Field(8), '^')
}

func (h OBXHelper) Probability() string {
	return h.seg.Field(9)
}

func (h OBXHelper) NatureOfAbnormalTest() string {
	return h.seg.Field(10)
}

func (h OBXHelper) ResultStatus() string {
	return h.seg.Field(11)
}

func (h OBXHelper) ObservationDateTime() string {
	return h.seg.Field(14)
}

func (h OBXHelper) ProducersID() string {
	return h.seg.Field(15)
}

func (h OBXHelper) ResponsibleObserver() string {
	return h.seg.Field(16)
}

func (h OBXHelper) ObservationMethod() string {
	return h.seg.Field(17)
}

func (m *Message) AllOBX() []OBXHelper {
	segments := m.Segments("OBX")
	result := make([]OBXHelper, len(segments))
	for i, seg := range segments {
		result[i] = OBXHelper{seg: seg}
	}
	return result
}

func (m *Message) AllNK1() []NK1Helper {
	segments := m.Segments("NK1")
	result := make([]NK1Helper, len(segments))
	for i, seg := range segments {
		result[i] = NK1Helper{seg: seg}
	}
	return result
}

func (m *Message) AllDG1() []DG1Helper {
	segments := m.Segments("DG1")
	result := make([]DG1Helper, len(segments))
	for i, seg := range segments {
		result[i] = DG1Helper{seg: seg}
	}
	return result
}

func (m *Message) AllOBR() []OBRHelper {
	segments := m.Segments("OBR")
	result := make([]OBRHelper, len(segments))
	for i, seg := range segments {
		result[i] = OBRHelper{seg: seg}
	}
	return result
}

type NK1Helper struct {
	seg Segment
}

func (m *Message) NK1() NK1Helper {
	if seg, ok := m.Segment("NK1"); ok {
		return NK1Helper{seg: seg}
	}
	return NK1Helper{}
}

func (h NK1Helper) Exists() bool {
	return h.seg.Name() != ""
}

func (h NK1Helper) SetID() string {
	return h.seg.Field(1)
}

func (h NK1Helper) Name() string {
	return h.seg.Field(2)
}

func (h NK1Helper) Relationship() string {
	return h.seg.Field(3)
}

func (h NK1Helper) Address() string {
	return h.seg.Field(4)
}

func (h NK1Helper) PhoneNumber() string {
	return h.seg.Field(5)
}

func (h NK1Helper) BusinessPhoneNumber() string {
	return h.seg.Field(6)
}

func (h NK1Helper) ContactRole() string {
	return h.seg.Field(7)
}

func (h NK1Helper) StartDate() string {
	return h.seg.Field(12)
}

func (h NK1Helper) EndDate() string {
	return h.seg.Field(13)
}

type DG1Helper struct {
	seg Segment
}

func (m *Message) DG1() DG1Helper {
	if seg, ok := m.Segment("DG1"); ok {
		return DG1Helper{seg: seg}
	}
	return DG1Helper{}
}

func (h DG1Helper) Exists() bool {
	return h.seg.Name() != ""
}

func (h DG1Helper) SetID() string {
	return h.seg.Field(1)
}

func (h DG1Helper) DiagnosisCode() string {
	return h.seg.Field(3)
}

func (h DG1Helper) DiagnosisDescription() string {
	return ParseComponent(h.seg.Field(4), 1)
}

func (h DG1Helper) DiagnosisType() string {
	return h.seg.Field(6)
}

func (h DG1Helper) DiagnosisDateTime() string {
	return h.seg.Field(5)
}
