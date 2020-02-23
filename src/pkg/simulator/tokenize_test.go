package simulator

import (
	"bufio"
	"github.com/apache/beam/sdks/go/pkg/beam"
	"github.com/apache/beam/sdks/go/pkg/beam/testing/ptest"
	"strings"
	"testing"
)

const (
	mockInput = `
DB00002	cetuximab	cetuximab,erbitux	-	retuximab	ritux	-	-	folfiri	-	-	colorectal_cancer_patients,hnscc,neoantigens,erbium
DB00003	dornase_alfa	pulmozyme	-	pulmozyne	-	-	-	-	-	-	hypersal,aerobika,aerochamber,airway_clearance,nebuliser
DB00005	etanercept	enbrel,etanercept	-	enbryl,ebrel,enbril,enbral,enbrell,enebrel,embrel	embril,evorel,labrel,lybrel	-	-	-	-	-	disease_modifying,immunomodulators,treat_ibd,treating_crohn's,antirheumatic,tnf_inhibitor,tnf_blockers,dmards,ra_patients,biologic_drugs,treat_multiple_sclerosis,biologic,immunomodulatory_drugs,tnf_inhibitors
DB00006	bivalirudin	angiomax,bivalirudin	-	-	arginmax	-	-	-	-	-	er_doc_wanted,wicksa,bp_monitoring,nurse_calls,angioma,crash_cart,partial_thromboplastin_time,distal_occlusion,pt_inr,iv_pole,post_cabg,outpatient_oncology,checking_blood_pressure,angiomas,phlebo,veinipuncture,peripherally_inserted_central_catheter,endotool,lab_draws,transvenous_pacemaker,nursing_assessments,12leads
DB00007	leuprolide	enantone,leuprolide,lupron_depot,microdose_lupron_protocol,leuprorelin,lupron,eligard,leuprolide_acetate	-	bupron,luprorelin,luprolide,leupron,lupon	lurpon,upfron,enantate,eligere,euphon	-	-	ivf_protocols	ivf_protocol,fet_protocol,microflare,ovarian_suppression	-	central_precocious_puberty,yale_biopsy,hystorectomy,person_stops_taking,hormone_replacement_drugs,12C50mg,remasculinize,rounds_of_iui,dominant_follicles,embryo_transfer,blockers,iui_or_ivf,frozen_embryo_transfers,therapist_is,antiandrogens,cycle_of_ivf,masculizing,fertility_testing,fertility_meds,ovarian_function,round_of_ivf,retrieval_cycle,injectable_form,ivf_attempt,depot_injection,superovulation,endo,natural_cycle_ivf,starting_injections,migard,six_cycles,estrogen_priming_protocol,upcoming_fet,simply_stop_taking,oversuppression,progestagens,18_mm_follicle,masculination,induce_ovulation,androgen_suppression,cycles_of_iui,3qd,replacement_therapy,injectible_iui,masculize,transbuccally,artificial_menopause,orichy,endometrial_receptivity_assay,hormonal_replacement,actual_fertility,dor_ladies,two_ivfs,mock_cycle,depot_medroxyprogesterone,estrogen_injections,fertility_test,fertility_doctor,fertility_drugs,freeze_some_sperm,blocking_drugs
`
)

func Test_tabbedLineHandler(t *testing.T) {
	buf := map[string]struct{}{}
	emit := func(token string) {
		buf[token] = struct{}{}
	}
	scanner := bufio.NewScanner(strings.NewReader(mockInput))
	for scanner.Scan() {
		tabbedLineHandler(scanner.Text(), emit)
	}

	for _, k := range []string{
		"cetuximab",
		"cetuximab,erbitux",
		"enbryl,ebrel,enbril,enbral,enbrell,enebrel,embrel",
	} {
		if _, ok := buf[k]; !ok {
			t.Errorf("%s expected but was not emitted", k)
		}
	}
}

func TestTokenize(t *testing.T) {
	want := map[string]struct{}{
		"cetuximab": struct{}{},
		"erbitux": struct{}{},
		"retuximab": struct{}{},
		"ritux": struct{}{},
		"folfiri": struct{}{},
		"erbium": struct{}{},
		"dornase_alfa": struct{}{},
		"pulmozyme": struct{}{},
		"hypersal": struct{}{},
		"aerobika": struct{}{},
		"etanercept": struct{}{},
		"enbrel": struct{}{},
		"enbril": struct{}{},
		"enbral": struct{}{},
		"enbrell": struct{}{},
		"enebrel": struct{}{},
		"embrel": struct{}{},
		"evorel": struct{}{},
		"labrel": struct{}{},
		"pulmozyne": struct {}{},
		"enbryl": struct {}{},
		"ebrel": struct {}{},
		"embril": struct {}{},
		"lybrel": struct {}{},
		"dmards": struct {}{},
		"bivalirudin": struct {}{},
		"angiomax": struct {}{},
		"epoetin": struct {}{},
		"darbepoetin_alfa": struct {}{},
		"leuprolide": struct {}{},
		"darbopoetin": struct {}{},
		"enantone": struct {}{},
		"lupron_depot": struct {}{},
		"microdose_lupron_protocol": struct {}{},
		"leuprorelin": struct {}{},
		"lupron": struct {}{},
		"eligard": struct {}{},
		"leuprolide_acetate": struct {}{},
		"bupron": struct {}{},
		"luprorelin": struct {}{},
		"luprolide": struct {}{},
		"leupron": struct {}{},
		"lupon": struct {}{},
		"-": struct {}{},
	}
	eval := func(got string, emit func(string)) {
		if _, ok := want[got]; !ok {
			t.Errorf("%s expected but not emitted", got)
		}
	}
	var input []interface{}
	scanner := bufio.NewScanner(strings.NewReader(mockInput))
	for scanner.Scan() {
		input = append(input, scanner.Text())
	}
	p, s, lines := ptest.Create(input)
	result := Tokenize(s, lines)
	beam.ParDo(s, eval, result)
	err := ptest.Run(p)
	if err != nil {
		t.Fatal(err)
	}
}