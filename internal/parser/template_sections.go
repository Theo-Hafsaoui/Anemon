package parser

const pro_item = "\\resumeItem{%ITEM%}\n"

const prof_template = `
\resumeSubheading
    {1}{2}
    {\href{3}{4}}{ }
\resumeItemListStart
    %ITEMS%
\resumeItemListEnd
`
const NB_P_PROF = 4 

const proj_template = `
\resumeProjectHeading
{\textbf{1} | \emph{2 \href{3}{\faIcon{github}}}}{}
\resumeItemListStart
    %ITEMS%
\resumeItemListEnd
`
const NB_P_PROJ = 3 

const edu_template = `
\resumeSubheading
{\href{2}{1}}{}
{3}{4}
`
const NB_P_EDU = 4 

const sk_template = `
\item \textbf{1}{: 2} \\
`
const NB_P_SK = 2 
