package knowledge

import "encoding/xml"

type CWEListImport struct {
	Weaknesses struct {
		XMLName    xml.Name      `xml:"Weaknesses"`
		Weaknesses []WeaknessCWE `xml:"Weakness"`
	} `xml:"Weaknesses"`
	Categories struct {
		XMLName    xml.Name   `xml:"Categories"`
		Categories []Category `xml:"Category"`
	} `xml:"Categories"`
}

type WeaknessCWE struct {
	XMLName             xml.Name `xml:"Weakness"`
	Text                string   `xml:",chardata"`
	ID                  string   `xml:"ID,attr"`
	Name                string   `xml:"Name,attr"`
	Abstraction         string   `xml:"Abstraction,attr"`
	Structure           string   `xml:"Structure,attr"`
	Status              string   `xml:"Status,attr"`
	Description         string   `xml:"Description"`
	MemberShips         []string
	ExtendedDescription struct {
		Text string   `xml:",chardata"`
		P    []string `xml:"p"`
		Ul   []struct {
			Text string `xml:",chardata"`
			Li   []struct {
				Text string `xml:",chardata"`
				Div  struct {
					Text string `xml:",chardata"`
					B    string `xml:"b"`
				} `xml:"div"`
				B string `xml:"b"`
			} `xml:"li"`
		} `xml:"ul"`
		Ol struct {
			Text string   `xml:",chardata"`
			Li   []string `xml:"li"`
		} `xml:"ol"`
		Br  []string `xml:"br"`
		Div struct {
			Text  string   `xml:",chardata"`
			Style string   `xml:"style,attr"`
			Div   []string `xml:"div"`
		} `xml:"div"`
	} `xml:"Extended_Description"`
	RelatedWeaknesses struct {
		Text            string `xml:",chardata"`
		RelatedWeakness []struct {
			Text    string `xml:",chardata"`
			Nature  string `xml:"Nature,attr"`
			CWEID   string `xml:"CWE_ID,attr"`
			ViewID  string `xml:"View_ID,attr"`
			Ordinal string `xml:"Ordinal,attr"`
			ChainID string `xml:"Chain_ID,attr"`
		} `xml:"Related_Weakness"`
	} `xml:"Related_Weaknesses"`
	ApplicablePlatforms struct {
		Text     string `xml:",chardata"`
		Language []struct {
			Text       string `xml:",chardata"`
			Class      string `xml:"Class,attr"`
			Prevalence string `xml:"Prevalence,attr"`
			Name       string `xml:"Name,attr"`
		} `xml:"Language"`
		Technology []struct {
			Text       string `xml:",chardata"`
			Class      string `xml:"Class,attr"`
			Prevalence string `xml:"Prevalence,attr"`
			Name       string `xml:"Name,attr"`
		} `xml:"Technology"`
		OperatingSystem []struct {
			Text       string `xml:",chardata"`
			Class      string `xml:"Class,attr"`
			Prevalence string `xml:"Prevalence,attr"`
			Name       string `xml:"Name,attr"`
		} `xml:"Operating_System"`
		Architecture []struct {
			Text       string `xml:",chardata"`
			Class      string `xml:"Class,attr"`
			Prevalence string `xml:"Prevalence,attr"`
			Name       string `xml:"Name,attr"`
		} `xml:"Architecture"`
	} `xml:"Applicable_Platforms"`
	BackgroundDetails struct {
		Text             string `xml:",chardata"`
		BackgroundDetail struct {
			Text string   `xml:",chardata"`
			P    []string `xml:"p"`
			Ul   struct {
				Text string   `xml:",chardata"`
				Li   []string `xml:"li"`
			} `xml:"ul"`
			Div struct {
				Text  string `xml:",chardata"`
				Style string `xml:"style,attr"`
				Br    string `xml:"br"`
			} `xml:"div"`
		} `xml:"Background_Detail"`
	} `xml:"Background_Details"`
	ModesOfIntroduction struct {
		Text         string `xml:",chardata"`
		Introduction []struct {
			Text  string `xml:",chardata"`
			Phase string `xml:"Phase"`
			Note  struct {
				Text string   `xml:",chardata"`
				P    []string `xml:"p"`
				Ul   []struct {
					Text string   `xml:",chardata"`
					Li   []string `xml:"li"`
				} `xml:"ul"`
			} `xml:"Note"`
		} `xml:"Introduction"`
	} `xml:"Modes_Of_Introduction"`
	LikelihoodOfExploit string `xml:"Likelihood_Of_Exploit"`
	CommonConsequences  struct {
		Text        string `xml:",chardata"`
		Consequence []struct {
			Text       string   `xml:",chardata"`
			Scope      []string `xml:"Scope"`
			Impact     []string `xml:"Impact"`
			Note       string   `xml:"Note"`
			Likelihood string   `xml:"Likelihood"`
		} `xml:"Consequence"`
	} `xml:"Common_Consequences"`
	DetectionMethods struct {
		Text            string `xml:",chardata"`
		DetectionMethod []struct {
			Text              string `xml:",chardata"`
			DetectionMethodID string `xml:"Detection_Method_ID,attr"`
			Method            string `xml:"Method"`
			Description       struct {
				Text string   `xml:",chardata"`
				P    []string `xml:"p"`
				Div  struct {
					Text  string   `xml:",chardata"`
					Style string   `xml:"style,attr"`
					Div   []string `xml:"div"`
					Ul    []struct {
						Text string   `xml:",chardata"`
						Li   []string `xml:"li"`
					} `xml:"ul"`
				} `xml:"div"`
				Ul []struct {
					Text string   `xml:",chardata"`
					Li   []string `xml:"li"`
				} `xml:"ul"`
			} `xml:"Description"`
			Effectiveness      string `xml:"Effectiveness"`
			EffectivenessNotes string `xml:"Effectiveness_Notes"`
		} `xml:"Detection_Method"`
	} `xml:"Detection_Methods"`
	PotentialMitigations struct {
		Text       string `xml:",chardata"`
		Mitigation []struct {
			Text         string   `xml:",chardata"`
			MitigationID string   `xml:"Mitigation_ID,attr"`
			Phase        []string `xml:"Phase"`
			Description  struct {
				Text string   `xml:",chardata"`
				P    []string `xml:"p"`
				Ul   []struct {
					Text string   `xml:",chardata"`
					Li   []string `xml:"li"`
				} `xml:"ul"`
				Div struct {
					Text  string   `xml:",chardata"`
					Style string   `xml:"style,attr"`
					Div   []string `xml:"div"`
				} `xml:"div"`
			} `xml:"Description"`
			Effectiveness      string `xml:"Effectiveness"`
			EffectivenessNotes struct {
				Text string `xml:",chardata"`
				P    string `xml:"p"`
			} `xml:"Effectiveness_Notes"`
			Strategy string `xml:"Strategy"`
		} `xml:"Mitigation"`
	} `xml:"Potential_Mitigations"`
	ObservedExamples struct {
		Text            string `xml:",chardata"`
		ObservedExample []struct {
			Text        string `xml:",chardata"`
			Reference   string `xml:"Reference"`
			Description string `xml:"Description"`
			Link        string `xml:"Link"`
		} `xml:"Observed_Example"`
	} `xml:"Observed_Examples"`
	References struct {
		Text      string `xml:",chardata"`
		Reference []struct {
			Text                string `xml:",chardata"`
			ExternalReferenceID string `xml:"External_Reference_ID,attr"`
			Section             string `xml:"Section,attr"`
		} `xml:"Reference"`
	} `xml:"References"`
	ContentHistory struct {
		Text       string `xml:",chardata"`
		Submission struct {
			Text                   string `xml:",chardata"`
			SubmissionName         string `xml:"Submission_Name"`
			SubmissionOrganization string `xml:"Submission_Organization"`
			SubmissionDate         string `xml:"Submission_Date"`
			SubmissionComment      string `xml:"Submission_Comment"`
		} `xml:"Submission"`
		Modification []struct {
			Text                     string `xml:",chardata"`
			ModificationName         string `xml:"Modification_Name"`
			ModificationOrganization string `xml:"Modification_Organization"`
			ModificationDate         string `xml:"Modification_Date"`
			ModificationComment      string `xml:"Modification_Comment"`
			ModificationImportance   string `xml:"Modification_Importance"`
		} `xml:"Modification"`
		PreviousEntryName []struct {
			Text string `xml:",chardata"`
			Date string `xml:"Date,attr"`
		} `xml:"Previous_Entry_Name"`
		Contribution []struct {
			Text                     string `xml:",chardata"`
			Type                     string `xml:"Type,attr"`
			ContributionName         string `xml:"Contribution_Name"`
			ContributionDate         string `xml:"Contribution_Date"`
			ContributionComment      string `xml:"Contribution_Comment"`
			ContributionOrganization string `xml:"Contribution_Organization"`
		} `xml:"Contribution"`
	} `xml:"Content_History"`
	WeaknessOrdinalities struct {
		Text               string `xml:",chardata"`
		WeaknessOrdinality []struct {
			Text        string `xml:",chardata"`
			Ordinality  string `xml:"Ordinality"`
			Description string `xml:"Description"`
		} `xml:"Weakness_Ordinality"`
	} `xml:"Weakness_Ordinalities"`
	AlternateTerms struct {
		Text          string `xml:",chardata"`
		AlternateTerm []struct {
			Text        string `xml:",chardata"`
			Term        string `xml:"Term"`
			Description struct {
				Text string   `xml:",chardata"`
				P    []string `xml:"p"`
				Ul   []struct {
					Text string   `xml:",chardata"`
					Li   []string `xml:"li"`
				} `xml:"ul"`
			} `xml:"Description"`
		} `xml:"Alternate_Term"`
	} `xml:"Alternate_Terms"`
	RelatedAttackPatterns struct {
		Text                 string `xml:",chardata"`
		RelatedAttackPattern []struct {
			Text    string `xml:",chardata"`
			CAPECID string `xml:"CAPEC_ID,attr"`
		} `xml:"Related_Attack_Pattern"`
	} `xml:"Related_Attack_Patterns"`
	TaxonomyMappings struct {
		Text            string `xml:",chardata"`
		TaxonomyMapping []struct {
			Text         string `xml:",chardata"`
			TaxonomyName string `xml:"Taxonomy_Name,attr"`
			EntryName    string `xml:"Entry_Name"`
			EntryID      string `xml:"Entry_ID"`
			MappingFit   string `xml:"Mapping_Fit"`
		} `xml:"Taxonomy_Mapping"`
	} `xml:"Taxonomy_Mappings"`
	Notes struct {
		Text string `xml:",chardata"`
		Note []struct {
			Text string   `xml:",chardata"`
			Type string   `xml:"Type,attr"`
			P    []string `xml:"p"`
			Ul   struct {
				Text string   `xml:",chardata"`
				Li   []string `xml:"li"`
			} `xml:"ul"`
			Div struct {
				Text  string   `xml:",chardata"`
				Style string   `xml:"style,attr"`
				Div   []string `xml:"div"`
			} `xml:"div"`
		} `xml:"Note"`
	} `xml:"Notes"`
	AffectedResources struct {
		Text             string   `xml:",chardata"`
		AffectedResource []string `xml:"Affected_Resource"`
	} `xml:"Affected_Resources"`
	FunctionalAreas struct {
		Text           string   `xml:",chardata"`
		FunctionalArea []string `xml:"Functional_Area"`
	} `xml:"Functional_Areas"`
}

type Category struct {
	XMLName xml.Name `xml:"Category"`
	Text    string   `xml:",chardata"`
	ID      string   `xml:"ID,attr"`
	Name    string   `xml:"Name,attr"`
	Status  string   `xml:"Status,attr"`
	Summary string   `xml:"Summary"`
	Notes   struct {
		Text string `xml:",chardata"`
		Note []struct {
			Text string   `xml:",chardata"`
			Type string   `xml:"Type,attr"`
			P    []string `xml:"p"`
			Ul   struct {
				Text string   `xml:",chardata"`
				Li   []string `xml:"li"`
			} `xml:"ul"`
		} `xml:"Note"`
	} `xml:"Notes"`
	ContentHistory struct {
		Text       string `xml:",chardata"`
		Submission struct {
			Text                   string `xml:",chardata"`
			SubmissionName         string `xml:"Submission_Name"`
			SubmissionDate         string `xml:"Submission_Date"`
			SubmissionComment      string `xml:"Submission_Comment"`
			SubmissionOrganization string `xml:"Submission_Organization"`
		} `xml:"Submission"`
		Modification []struct {
			Text                     string `xml:",chardata"`
			ModificationName         string `xml:"Modification_Name"`
			ModificationOrganization string `xml:"Modification_Organization"`
			ModificationDate         string `xml:"Modification_Date"`
			ModificationComment      string `xml:"Modification_Comment"`
		} `xml:"Modification"`
		PreviousEntryName []struct {
			Text string `xml:",chardata"`
			Date string `xml:"Date,attr"`
		} `xml:"Previous_Entry_Name"`
		Contribution []struct {
			Text                     string `xml:",chardata"`
			Type                     string `xml:"Type,attr"`
			ContributionName         string `xml:"Contribution_Name"`
			ContributionOrganization string `xml:"Contribution_Organization"`
			ContributionDate         string `xml:"Contribution_Date"`
			ContributionComment      string `xml:"Contribution_Comment"`
		} `xml:"Contribution"`
	} `xml:"Content_History"`
	Relationships struct {
		Text      string `xml:",chardata"`
		HasMember []struct {
			Text   string `xml:",chardata"`
			CWEID  string `xml:"CWE_ID,attr"`
			ViewID string `xml:"View_ID,attr"`
		} `xml:"Has_Member"`
	} `xml:"Relationships"`
	References struct {
		Text      string `xml:",chardata"`
		Reference []struct {
			Text                string `xml:",chardata"`
			ExternalReferenceID string `xml:"External_Reference_ID,attr"`
			Section             string `xml:"Section,attr"`
		} `xml:"Reference"`
	} `xml:"References"`
	TaxonomyMappings struct {
		Text            string `xml:",chardata"`
		TaxonomyMapping []struct {
			Text         string `xml:",chardata"`
			TaxonomyName string `xml:"Taxonomy_Name,attr"`
			EntryID      string `xml:"Entry_ID"`
			EntryName    string `xml:"Entry_Name"`
			MappingFit   string `xml:"Mapping_Fit"`
		} `xml:"Taxonomy_Mapping"`
	} `xml:"Taxonomy_Mappings"`
}
