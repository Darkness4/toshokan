package opds_test

import (
	"encoding/xml"
	"fmt"
	"testing"

	"github.com/Darkness4/toshokan/opds"
	"github.com/stretchr/testify/assert"
)

const fixtureEntryFeed = `<entry xmlns="http://www.w3.org/2005/Atom" xmlns:thr="http://purl.org/syndication/thread/1.0" xmlns:dcterms="http://purl.org/dc/terms/" xmlns:opds="http://opds-spec.org/2010/catalog" xmlns:pse="http://vaemendis.net/opds-pse/ns" xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance" xmlns:schema="http://schema.org/">
<link rel="start" href="/api/opds" type="application/atom+xml;profile=opds-catalog;kind=navigation"/>
<link rel="self" href="/api/opds/9fb486548e8a248a06438a9d62ddbc1bc333ada0" type="application/atom+xml;type=entry;profile=opds-catalog"/>
<title>24 Hour Humiliation TV</title>
<id>urn:lrr:9fb486548e8a248a06438a9d62ddbc1bc333ada0</id>
<updated>2024-05-26T20:37:36Z</updated>
<published>2024-05-26T20:37:36Z</published>
<author>
<name>Crimson</name>
</author>
<rights/>
<dcterms:language>english</dcterms:language>
<dcterms:publisher/>
<dcterms:issued/>
<category term="New Archive"></category>
<summary>tags</summary>
<link rel="http://opds-spec.org/image" href="/api/archives/9fb486548e8a248a06438a9d62ddbc1bc333ada0/thumbnail" type="image/jpeg"/>
<link rel="http://opds-spec.org/image/thumbnail" href="/api/archives/9fb486548e8a248a06438a9d62ddbc1bc333ada0/thumbnail" type="image/jpeg"/>
<link rel="http://opds-spec.org/acquisition" href="/api/archives/9fb486548e8a248a06438a9d62ddbc1bc333ada0/download" title="Download/Read" type="application/x-cbz"/>
<link rel="http://vaemendis.net/opds-pse/stream" type="image/jpeg" href="/api/opds/9fb486548e8a248a06438a9d62ddbc1bc333ada0/pse?page={pageNumber}" pse:count="78"/>
<link type="text/html" rel="alternate" title="Open in LANraragi" href="/reader?id=9fb486548e8a248a06438a9d62ddbc1bc333ada0"/>
</entry>`

func TestParseEntryXML(t *testing.T) {
	e1 := opds.Entry{}
	err := xml.Unmarshal([]byte(fixtureEntryFeed), &e1)

	assert.NoError(t, err)

	b, err := xml.MarshalIndent(e1, "", "  ")
	assert.NoError(t, err)
	fmt.Println(string(b))

	e2 := opds.Entry{}
	err = xml.Unmarshal(b, &e2)
	assert.NoError(t, err)
	assert.Equal(t, e1, e2)
}
