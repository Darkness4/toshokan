package opds

import (
	"encoding/xml"

	"github.com/gorilla/feeds"
)

const thrns = "http://purl.org/syndication/thread/1.0"
const dctermsns = "http://purl.org/dc/terms/"
const opdsns = "http://opds-spec.org/2010/catalog"
const psens = "http://vaemendis.net/opds-pse/ns"
const xsins = "http://www.w3.org/2001/XMLSchema-instance"
const schemans = "http://schema.org/"
const ns = "http://www.w3.org/2005/Atom"

type Link struct {
	XMLName xml.Name `xml:"link"`
	Href    string   `xml:"href,attr"`
	Rel     string   `xml:"rel,attr,omitempty"`
	Type    string   `xml:"type,attr,omitempty"`
	Length  string   `xml:"length,attr,omitempty"`
	Title   string   `xml:"title,attr,omitempty"`

	PSECount        string `xml:"pse:count,attr,omitempty"`
	PSELastRead     string `xml:"pse:lastRead,attr,omitempty"`
	PSELastReadDate string `xml:"pse:lastReadDate,attr,omitempty"`

	OPDSFacetGroup  string `xml:"opds:facetGroup,attr,omitempty"`
	OPDSActiveFacet string `xml:"opds:activeFacet,attr,omitempty"`

	ThrCount string `xml:"thr:count,attr,omitempty"`
}

type Category struct {
	XMLName xml.Name `xml:"category"`
	Term    string   `xml:"term,attr"`
	Content string   `xml:",chardata"`
}

type Summary struct {
	XMLName xml.Name `xml:"summary"`
	Content string   `xml:",chardata"`
	Type    string   `xml:"type,attr,omitempty"`
}

type Entry struct {
	XMLName      xml.Name `xml:"entry"`
	Xmlns        string   `xml:"xmlns,attr,omitempty"`
	ThrXmlns     string   `xml:"xmlns:thr,attr,omitempty"`
	DctermsXmlns string   `xml:"xmlns:dcterms,attr,omitempty"`
	OpdsXmlns    string   `xml:"xmlns:opds,attr,omitempty"`
	PseXmlns     string   `xml:"xmlns:pse,attr,omitempty"`
	XsiXmlns     string   `xml:"xmlns:xsi,attr,omitempty"`
	SchemaXmlns  string   `xml:"xmlns:schema,attr,omitempty"`

	Title       string   `xml:"title"`
	Updated     string   `xml:"updated"`
	ID          string   `xml:"id"`
	Category    Category `xml:"category,omitempty"`
	Content     *feeds.AtomContent
	Rights      string `xml:"rights"`
	Source      string `xml:"source,omitempty"`
	Published   string `xml:"published,omitempty"`
	Contributor *feeds.AtomContributor

	Summary   *Summary          // required if content has src or content is base64
	Author    *feeds.AtomAuthor // required if feed lacks an author
	Language  string            `xml:"dcterms:language"`
	Publisher string            `xml:"dcterms:publisher"`
	Issued    string            `xml:"dcterms:issued"`

	Links []Link `xml:"link"`
}

func (e *Entry) Enveloppe() *Entry {
	e.ThrXmlns = thrns
	e.DctermsXmlns = dctermsns
	e.OpdsXmlns = opdsns
	e.PseXmlns = psens
	e.XsiXmlns = xsins
	e.Xmlns = ns
	e.SchemaXmlns = schemans
	return e
}

type Feed struct {
	XMLName      xml.Name `xml:"feed"`
	Xmlns        string   `xml:"xmlns,attr"`
	ThrXmlns     string   `xml:"xmlns:thr,attr"`
	DctermsXmlns string   `xml:"xmlns:dcterms,attr"`
	OpdsXmlns    string   `xml:"xmlns:opds,attr"`
	PseXmlns     string   `xml:"xmlns:pse,attr"`

	Title    string `xml:"title"`   // required
	ID       string `xml:"id"`      // required
	Updated  string `xml:"updated"` // required
	Category string `xml:"category,omitempty"`
	Icon     string `xml:"icon,omitempty"`
	Logo     string `xml:"logo,omitempty"`
	Rights   string `xml:"rights,omitempty"` // copyright used
	Subtitle string `xml:"subtitle,omitempty"`

	Author      *feeds.AtomAuthor `xml:"author,omitempty"`
	Contributor *feeds.AtomContributor

	Links []Link `xml:"link"`
}

func (f *Feed) Enveloppe() *Feed {
	f.ThrXmlns = thrns
	f.DctermsXmlns = dctermsns
	f.OpdsXmlns = opdsns
	f.PseXmlns = psens
	f.Xmlns = ns
	return f
}
