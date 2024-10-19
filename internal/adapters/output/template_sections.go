package output

type TemplateStruct struct {
	template string
	hook     string
}

var (
	single_item_template = "\\resumeItem{%ITEM%}\n"

	ProfessionalTemplate = TemplateStruct{
		template: `
\resumeSubheading
    {$1}{$2}
    {\href{$3}{$4}}{ }
\resumeItemListStart
    %ITEMS%
\resumeItemListEnd
`,
		hook: "%EXPERIENCE_SECTIONS%",
	}

	ProjectTemplate = TemplateStruct{
		template: `
\resumeProjectHeading
{\textbf{$1} | \emph{$2 \href{$3}{\faIcon{github}}}}{}
\resumeItemListStart
    %ITEMS%
\resumeItemListEnd
`,
		hook: "%PROJECTS_SECTIONS%",
	}

	EducationTemplate = TemplateStruct{
		template: `
\resumeSubheading
{\href{$2}{$1}}{}
{$3}{$4}
`,
		hook: "%EDUCATION_SECTIONS%",
	}

	SkillTemplate = TemplateStruct{
		template: `
\item \textbf{$1}{: $2} \\
`,
		hook: "%SKILL_SECTIONS%",
	}
)

// For section translation
var translations = map[string]map[string]string{
    "FR": {
        "Skill":       "Compétence",
        "Education":   "Éducation",
        "Professional": "Professionnel",
        "Project":     "Projet",
    },
    "DE": {
        "Skill":       "Fähigkeit",
        "Education":   "Bildung",
        "Professional": "Beruflich",
        "Project":     "Projekt",
    },
    "NL": {
        "Skill":       "Vaardigheid",
        "Education":   "Onderwijs",
        "Professional": "Professioneel",
        "Project":     "Project",
    },
    "SP": {
        "Skill":       "Habilidad",
        "Education":   "Educación",
        "Professional": "Profesional",
        "Project":     "Proyecto",
    },
    "IT": {
        "Skill":       "Competenza",
        "Education":   "Istruzione",
        "Professional": "Professionale",
        "Project":     "Progetto",
    },
    "PT": {
        "Skill":       "Habilidade",
        "Education":   "Educação",
        "Professional": "Profissional",
        "Project":     "Projeto",
    },
    "SE": {
        "Skill":       "Färdighet",
        "Education":   "Utbildning",
        "Professional": "Professionell",
        "Project":     "Projekt",
    },
    "NO": {
        "Skill":       "Ferdighet",
        "Education":   "Utdanning",
        "Professional": "Profesjonell",
        "Project":     "Prosjekt",
    },
    "FI": {
        "Skill":       "Taito",
        "Education":   "Koulutus",
        "Professional": "Ammattilainen",
        "Project":     "Hanke",
    },
    "PL": {
        "Skill":       "Umiejętność",
        "Education":   "Edukacja",
        "Professional": "Profesjonalny",
        "Project":     "Projekt",
    },
    "RU": {
        "Skill":       "Навыки",
        "Education":   "Образование",
        "Professional": "Профессионал",
        "Project":     "Проект",
    },
}
