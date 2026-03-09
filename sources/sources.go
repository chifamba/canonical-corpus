// Package sources defines the full catalogue of biblical corpus books and their
// public-domain source URLs.
package sources

import (
	"fmt"
	"net/url"

	"github.com/chifamba/canonical-corpus/internal/metadata"
)

// ---------------------------------------------------------------------------
// Source-reference helpers
// ---------------------------------------------------------------------------

func wikisource(path string) metadata.SourceRef {
	return metadata.SourceRef{
		URL:      "https://en.wikisource.org/wiki/" + path,
		Format:   "html",
		Language: "en",
		License:  "Public Domain",
	}
}

func sacredTexts(path string) metadata.SourceRef {
	return metadata.SourceRef{
		URL:      "https://www.sacred-texts.com/" + path,
		Format:   "html",
		Language: "en",
		License:  "Public Domain",
	}
}

func gutenberg(id string) metadata.SourceRef {
	return metadata.SourceRef{
		URL:      "https://www.gutenberg.org/files/" + id,
		Format:   "txt",
		Language: "en",
		License:  "Public Domain",
	}
}

// webSource returns a World English Bible (WEB) source from sacred-texts.com.
// abbr is the lowercase OSIS-style book abbreviation used by that site (e.g. "gen", "mat").
func webSource(abbr string) metadata.SourceRef {
	return metadata.SourceRef{
		URL:        "https://sacred-texts.com/bib/web/" + abbr + ".htm",
		Format:     "html",
		Language:   "en",
		License:    "Public Domain",
		Translator: "Rainbow Missions, Inc.",
		Notes:      "World English Bible (WEB) — modern public-domain English, closest to NIV",
	}
}

// asvSource returns an American Standard Version (ASV) source from Wikisource.
// wsPath is the page name under Bible_(American_Standard)/ (e.g. "Genesis", "I_Samuel").
func asvSource(wsPath string) metadata.SourceRef {
	return metadata.SourceRef{
		URL:        "https://en.wikisource.org/wiki/Bible_(American_Standard)/" + wsPath,
		Format:     "html",
		Language:   "en",
		License:    "Public Domain",
		Translator: "American Standard Version Revision Committee",
		Notes:      "American Standard Version (ASV, 1901) — formal-equivalence, closest to ESV",
	}
}

// lxxSource returns a Greek Septuagint source from sacred-texts.com Polyglot Bible.
// abbr is the book abbreviation used by the polyglot section (e.g. "gen", "mat").
func lxxSource(abbr string) metadata.SourceRef {
	return metadata.SourceRef{
		URL:      "https://www.sacred-texts.com/bib/poly/" + abbr + ".htm",
		Format:   "html",
		Language: "el",
		License:  "Public Domain",
		Notes:    "Polyglot Bible — Greek Septuagint (LXX, OT) / Textus Receptus (NT) columns",
	}
}

// mechonMamre returns a Hebrew Masoretic Text source from Mechon Mamre.
// code is the two-digit book number used by that site (e.g. "01" for Genesis).
func mechonMamre(code string) metadata.SourceRef {
	return metadata.SourceRef{
		URL:      "https://mechon-mamre.org/p/pt/pt" + code + ".htm",
		Format:   "html",
		Language: "he",
		License:  "Public Domain",
		Notes:    "Westminster Leningrad Codex via Mechon Mamre — modern-vowelled Hebrew Masoretic Text",
	}
}

// shonaSource returns a Shona Bible source from eBible.org.
// abbr is the book abbreviation used by that site (e.g. "GEN", "MAT").
func shonaSource(abbr string) metadata.SourceRef {
	return metadata.SourceRef{
		URL:      "https://ebible.org/snb/" + abbr + ".htm",
		Format:   "html",
		Language: "sn",
		License:  "Public Domain",
		Notes:    "Shona New Bible (SNB) — public-domain Shona translation",
	}
}

func superSearchSource(module, bookName string) metadata.SourceRef {
	return metadata.SourceRef{
		URL:     fmt.Sprintf("https://api.biblesupersearch.com/api?bible=%s&reference=%s", module, url.QueryEscape(bookName)),
		Format:  "json",
		License: "Public Domain / Freely Licensed",
		Notes:   "Fetched via Bible SuperSearch API",
	}
}

// ---------------------------------------------------------------------------
// Book-entry builders  (one per translation)
// ---------------------------------------------------------------------------

func genericSuperSearchBook(order int, title, module, lang, notes string) metadata.BookMeta {
	return metadata.BookMeta{
		Title:          title,
		CanonicalOrder: order,
		Category:       metadata.CategoryCanonical,
		Language:       lang,
		TranslationID:  module,
		License:        "Public Domain / Freely Licensed",
		Notes:          notes,
		Sources:        []metadata.SourceRef{superSearchSource(module, title)},
	}
}

// yltBook returns a Young's Literal Translation book.
func yltBook(order int, title string) metadata.BookMeta {
	return genericSuperSearchBook(order, title, "ylt", "en", "Young's Literal Translation (1862/1898)")
}

// dbtBook returns a Darby Bible book.
func dbtBook(order int, title string) metadata.BookMeta {
	return genericSuperSearchBook(order, title, "dbt", "en", "Darby Bible (1890)")
}

// drbBook returns a Douay-Rheims Bible book.
func drbBook(order int, title string) metadata.BookMeta {
	return genericSuperSearchBook(order, title, "drb", "en", "Douay-Rheims Bible (1899)")
}

// trBook returns a Textus Receptus (Greek) book.
func trBook(order int, title string) metadata.BookMeta {
	return genericSuperSearchBook(order, title, "tr", "el", "Textus Receptus Greek New Testament")
}

// kjvBook returns a canonical KJV entry sourced from Wikisource.
func kjvBook(order int, title, wsPath string, extra ...metadata.SourceRef) metadata.BookMeta {
	srcs := []metadata.SourceRef{wikisource(wsPath)}
	srcs = append(srcs, extra...)
	return metadata.BookMeta{
		Title:          title,
		CanonicalOrder: order,
		Category:       metadata.CategoryCanonical,
		Language:       "en",
		TranslationID:  "kjv",
		License:        "Public Domain",
		Translator:     "King James Version Translators",
		Sources:        srcs,
	}
}

// webBook returns a canonical WEB entry from sacred-texts.com.
func webBook(order int, title, abbr string) metadata.BookMeta {
	return metadata.BookMeta{
		Title:          title,
		CanonicalOrder: order,
		Category:       metadata.CategoryCanonical,
		Language:       "en",
		TranslationID:  "web",
		License:        "Public Domain",
		Translator:     "Rainbow Missions, Inc.",
		Sources:        []metadata.SourceRef{webSource(abbr)},
	}
}

// asvBook returns a canonical ASV entry from Wikisource.
func asvBook(order int, title, wsPath string) metadata.BookMeta {
	return metadata.BookMeta{
		Title:          title,
		CanonicalOrder: order,
		Category:       metadata.CategoryCanonical,
		Language:       "en",
		TranslationID:  "asv",
		License:        "Public Domain",
		Translator:     "American Standard Version Revision Committee",
		Sources:        []metadata.SourceRef{asvSource(wsPath)},
	}
}

// greekBook returns a canonical Greek (LXX/GNT) entry from sacred-texts.com Polyglot.
func greekBook(order int, title, polyAbbr string) metadata.BookMeta {
	return metadata.BookMeta{
		Title:          title,
		CanonicalOrder: order,
		Category:       metadata.CategoryCanonical,
		Language:       "el",
		TranslationID:  "lxx",
		License:        "Public Domain",
		Notes:          "Greek Septuagint (OT) / Greek New Testament (NT) — via Polyglot Bible",
		Sources:        []metadata.SourceRef{lxxSource(polyAbbr)},
	}
}

// hebrewBook returns a canonical Hebrew Masoretic Text entry from Mechon Mamre.
func hebrewBook(order int, title, mmCode string) metadata.BookMeta {
	return metadata.BookMeta{
		Title:          title,
		CanonicalOrder: order,
		Category:       metadata.CategoryCanonical,
		Language:       "he",
		TranslationID:  "hmt",
		License:        "Public Domain",
		Notes:          "Hebrew Masoretic Text — Westminster Leningrad Codex via Mechon Mamre",
		Sources:        []metadata.SourceRef{mechonMamre(mmCode)},
	}
}

// wlcBook returns a canonical Hebrew Masoretic Text entry representing the WLC.
func wlcBook(order int, title, mmCode string) metadata.BookMeta {
	return metadata.BookMeta{
		Title:          title,
		CanonicalOrder: order,
		Category:       metadata.CategoryCanonical,
		Language:       "he",
		TranslationID:  "wlc",
		License:        "Public Domain",
		Notes:          "Westminster Leningrad Codex (WLC) via Mechon Mamre",
		Sources:        []metadata.SourceRef{mechonMamre(mmCode)},
	}
}

func sblgSource(filename string) metadata.SourceRef {
	return metadata.SourceRef{
		URL:      "https://raw.githubusercontent.com/Faithlife/SBLGNT/master/data/sblgnt/text/" + filename,
		Format:   "txt",
		Language: "el",
		License:  "CC-BY 4.0",
		Notes:    "SBL Greek New Testament",
	}
}

// sblgBook returns a canonical SBL Greek New Testament entry.
func sblgBook(order int, title, filename string) metadata.BookMeta {
	return metadata.BookMeta{
		Title:          title,
		CanonicalOrder: order,
		Category:       metadata.CategoryCanonical,
		Language:       "el",
		TranslationID:  "sblg",
		License:        "CC-BY 4.0",
		Translator:     "Michael W. Holmes (SBL/Logos)",
		Notes:          "SBL Greek New Testament",
		Sources:        []metadata.SourceRef{sblgSource(filename)},
	}
}

// shonaBook returns a canonical Shona Bible entry from eBible.org.
func shonaBook(order int, title, abbr string) metadata.BookMeta {
	return metadata.BookMeta{
		Title:          title,
		CanonicalOrder: order,
		Category:       metadata.CategoryCanonical,
		Language:       "sn",
		TranslationID:  "sna",
		License:        "CC BY-SA 4.0",
		Translator:     "Biblica, Inc.",
		Notes:          "Biblica® Bhaibheri Dzvene Rakasununguka MuChiShona Chanhasi 2017",
		Sources:        []metadata.SourceRef{shonaSource(abbr)},
	}
}

// ---------------------------------------------------------------------------
// Canonical Books — KJV (Protestant 66-book canon)
// ---------------------------------------------------------------------------

var kjvCanonicalBooks = []metadata.BookMeta{
	// ---- Old Testament ----
	kjvBook(1, "Genesis", "Bible_(King_James)/Genesis"),
	kjvBook(2, "Exodus", "Bible_(King_James)/Exodus"),
	kjvBook(3, "Leviticus", "Bible_(King_James)/Leviticus"),
	kjvBook(4, "Numbers", "Bible_(King_James)/Numbers"),
	kjvBook(5, "Deuteronomy", "Bible_(King_James)/Deuteronomy"),
	kjvBook(6, "Joshua", "Bible_(King_James)/Joshua"),
	kjvBook(7, "Judges", "Bible_(King_James)/Judges"),
	kjvBook(8, "Ruth", "Bible_(King_James)/Ruth"),
	kjvBook(9, "1 Samuel", "Bible_(King_James)/I_Samuel"),
	kjvBook(10, "2 Samuel", "Bible_(King_James)/II_Samuel"),
	kjvBook(11, "1 Kings", "Bible_(King_James)/I_Kings"),
	kjvBook(12, "2 Kings", "Bible_(King_James)/II_Kings"),
	kjvBook(13, "1 Chronicles", "Bible_(King_James)/I_Chronicles"),
	kjvBook(14, "2 Chronicles", "Bible_(King_James)/II_Chronicles"),
	kjvBook(15, "Ezra", "Bible_(King_James)/Ezra"),
	kjvBook(16, "Nehemiah", "Bible_(King_James)/Nehemiah"),
	kjvBook(17, "Esther", "Bible_(King_James)/Esther"),
	kjvBook(18, "Job", "Bible_(King_James)/Job"),
	kjvBook(19, "Psalms", "Bible_(King_James)/Psalms"),
	kjvBook(20, "Proverbs", "Bible_(King_James)/Proverbs"),
	kjvBook(21, "Ecclesiastes", "Bible_(King_James)/Ecclesiastes"),
	kjvBook(22, "Song of Solomon", "Bible_(King_James)/Song_of_Solomon"),
	kjvBook(23, "Isaiah", "Bible_(King_James)/Isaiah"),
	kjvBook(24, "Jeremiah", "Bible_(King_James)/Jeremiah"),
	kjvBook(25, "Lamentations", "Bible_(King_James)/Lamentations"),
	kjvBook(26, "Ezekiel", "Bible_(King_James)/Ezekiel"),
	kjvBook(27, "Daniel", "Bible_(King_James)/Daniel"),
	kjvBook(28, "Hosea", "Bible_(King_James)/Hosea"),
	kjvBook(29, "Joel", "Bible_(King_James)/Joel"),
	kjvBook(30, "Amos", "Bible_(King_James)/Amos"),
	kjvBook(31, "Obadiah", "Bible_(King_James)/Obadiah"),
	kjvBook(32, "Jonah", "Bible_(King_James)/Jonah"),
	kjvBook(33, "Micah", "Bible_(King_James)/Micah"),
	kjvBook(34, "Nahum", "Bible_(King_James)/Nahum"),
	kjvBook(35, "Habakkuk", "Bible_(King_James)/Habakkuk"),
	kjvBook(36, "Zephaniah", "Bible_(King_James)/Zephaniah"),
	kjvBook(37, "Haggai", "Bible_(King_James)/Haggai"),
	kjvBook(38, "Zechariah", "Bible_(King_James)/Zechariah"),
	kjvBook(39, "Malachi", "Bible_(King_James)/Malachi"),
	// ---- New Testament ----
	kjvBook(40, "Matthew", "Bible_(King_James)/Matthew"),
	kjvBook(41, "Mark", "Bible_(King_James)/Mark"),
	kjvBook(42, "Luke", "Bible_(King_James)/Luke"),
	kjvBook(43, "John", "Bible_(King_James)/John"),
	kjvBook(44, "Acts", "Bible_(King_James)/Acts"),
	kjvBook(45, "Romans", "Bible_(King_James)/Romans"),
	kjvBook(46, "1 Corinthians", "Bible_(King_James)/I_Corinthians"),
	kjvBook(47, "2 Corinthians", "Bible_(King_James)/II_Corinthians"),
	kjvBook(48, "Galatians", "Bible_(King_James)/Galatians"),
	kjvBook(49, "Ephesians", "Bible_(King_James)/Ephesians"),
	kjvBook(50, "Philippians", "Bible_(King_James)/Philippians"),
	kjvBook(51, "Colossians", "Bible_(King_James)/Colossians"),
	kjvBook(52, "1 Thessalonians", "Bible_(King_James)/I_Thessalonians"),
	kjvBook(53, "2 Thessalonians", "Bible_(King_James)/II_Thessalonians"),
	kjvBook(54, "1 Timothy", "Bible_(King_James)/I_Timothy"),
	kjvBook(55, "2 Timothy", "Bible_(King_James)/II_Timothy"),
	kjvBook(56, "Titus", "Bible_(King_James)/Titus"),
	kjvBook(57, "Philemon", "Bible_(King_James)/Philemon"),
	kjvBook(58, "Hebrews", "Bible_(King_James)/Hebrews"),
	kjvBook(59, "James", "Bible_(King_James)/James"),
	kjvBook(60, "1 Peter", "Bible_(King_James)/I_Peter"),
	kjvBook(61, "2 Peter", "Bible_(King_James)/II_Peter"),
	kjvBook(62, "1 John", "Bible_(King_James)/I_John"),
	kjvBook(63, "2 John", "Bible_(King_James)/II_John"),
	kjvBook(64, "3 John", "Bible_(King_James)/III_John"),
	kjvBook(65, "Jude", "Bible_(King_James)/Jude"),
	kjvBook(66, "Revelation", "Bible_(King_James)/Revelation"),
}

// ---------------------------------------------------------------------------
// Canonical Books — WEB (World English Bible, ~NIV equivalent)
// ---------------------------------------------------------------------------

var webCanonicalBooks = []metadata.BookMeta{
	// ---- Old Testament ----
	webBook(1, "Genesis", "gen"),
	webBook(2, "Exodus", "exo"),
	webBook(3, "Leviticus", "lev"),
	webBook(4, "Numbers", "num"),
	webBook(5, "Deuteronomy", "deu"),
	webBook(6, "Joshua", "jos"),
	webBook(7, "Judges", "jdg"),
	webBook(8, "Ruth", "rut"),
	webBook(9, "1 Samuel", "1sa"),
	webBook(10, "2 Samuel", "2sa"),
	webBook(11, "1 Kings", "1ki"),
	webBook(12, "2 Kings", "2ki"),
	webBook(13, "1 Chronicles", "1ch"),
	webBook(14, "2 Chronicles", "2ch"),
	webBook(15, "Ezra", "ezr"),
	webBook(16, "Nehemiah", "neh"),
	webBook(17, "Esther", "est"),
	webBook(18, "Job", "job"),
	webBook(19, "Psalms", "psa"),
	webBook(20, "Proverbs", "pro"),
	webBook(21, "Ecclesiastes", "ecc"),
	webBook(22, "Song of Solomon", "sng"),
	webBook(23, "Isaiah", "isa"),
	webBook(24, "Jeremiah", "jer"),
	webBook(25, "Lamentations", "lam"),
	webBook(26, "Ezekiel", "ezk"),
	webBook(27, "Daniel", "dan"),
	webBook(28, "Hosea", "hos"),
	webBook(29, "Joel", "joe"),
	webBook(30, "Amos", "amo"),
	webBook(31, "Obadiah", "oba"),
	webBook(32, "Jonah", "jon"),
	webBook(33, "Micah", "mic"),
	webBook(34, "Nahum", "nah"),
	webBook(35, "Habakkuk", "hab"),
	webBook(36, "Zephaniah", "zep"),
	webBook(37, "Haggai", "hag"),
	webBook(38, "Zechariah", "zec"),
	webBook(39, "Malachi", "mal"),
	// ---- New Testament ----
	webBook(40, "Matthew", "mat"),
	webBook(41, "Mark", "mrk"),
	webBook(42, "Luke", "luk"),
	webBook(43, "John", "jhn"),
	webBook(44, "Acts", "act"),
	webBook(45, "Romans", "rom"),
	webBook(46, "1 Corinthians", "1co"),
	webBook(47, "2 Corinthians", "2co"),
	webBook(48, "Galatians", "gal"),
	webBook(49, "Ephesians", "eph"),
	webBook(50, "Philippians", "php"),
	webBook(51, "Colossians", "col"),
	webBook(52, "1 Thessalonians", "1th"),
	webBook(53, "2 Thessalonians", "2th"),
	webBook(54, "1 Timothy", "1ti"),
	webBook(55, "2 Timothy", "2ti"),
	webBook(56, "Titus", "tit"),
	webBook(57, "Philemon", "phm"),
	webBook(58, "Hebrews", "heb"),
	webBook(59, "James", "jas"),
	webBook(60, "1 Peter", "1pe"),
	webBook(61, "2 Peter", "2pe"),
	webBook(62, "1 John", "1jo"),
	webBook(63, "2 John", "2jo"),
	webBook(64, "3 John", "3jo"),
	webBook(65, "Jude", "jud"),
	webBook(66, "Revelation", "rev"),
}

// ---------------------------------------------------------------------------
// Canonical Books — ASV (American Standard Version, ~ESV equivalent)
// ---------------------------------------------------------------------------

var asvCanonicalBooks = []metadata.BookMeta{
	// ---- Old Testament ----
	asvBook(1, "Genesis", "Genesis"),
	asvBook(2, "Exodus", "Exodus"),
	asvBook(3, "Leviticus", "Leviticus"),
	asvBook(4, "Numbers", "Numbers"),
	asvBook(5, "Deuteronomy", "Deuteronomy"),
	asvBook(6, "Joshua", "Joshua"),
	asvBook(7, "Judges", "Judges"),
	asvBook(8, "Ruth", "Ruth"),
	asvBook(9, "1 Samuel", "I_Samuel"),
	asvBook(10, "2 Samuel", "II_Samuel"),
	asvBook(11, "1 Kings", "I_Kings"),
	asvBook(12, "2 Kings", "II_Kings"),
	asvBook(13, "1 Chronicles", "I_Chronicles"),
	asvBook(14, "2 Chronicles", "II_Chronicles"),
	asvBook(15, "Ezra", "Ezra"),
	asvBook(16, "Nehemiah", "Nehemiah"),
	asvBook(17, "Esther", "Esther"),
	asvBook(18, "Job", "Job"),
	asvBook(19, "Psalms", "Psalms"),
	asvBook(20, "Proverbs", "Proverbs"),
	asvBook(21, "Ecclesiastes", "Ecclesiastes"),
	asvBook(22, "Song of Solomon", "Song_of_Solomon"),
	asvBook(23, "Isaiah", "Isaiah"),
	asvBook(24, "Jeremiah", "Jeremiah"),
	asvBook(25, "Lamentations", "Lamentations"),
	asvBook(26, "Ezekiel", "Ezekiel"),
	asvBook(27, "Daniel", "Daniel"),
	asvBook(28, "Hosea", "Hosea"),
	asvBook(29, "Joel", "Joel"),
	asvBook(30, "Amos", "Amos"),
	asvBook(31, "Obadiah", "Obadiah"),
	asvBook(32, "Jonah", "Jonah"),
	asvBook(33, "Micah", "Micah"),
	asvBook(34, "Nahum", "Nahum"),
	asvBook(35, "Habakkuk", "Habakkuk"),
	asvBook(36, "Zephaniah", "Zephaniah"),
	asvBook(37, "Haggai", "Haggai"),
	asvBook(38, "Zechariah", "Zechariah"),
	asvBook(39, "Malachi", "Malachi"),
	// ---- New Testament ----
	asvBook(40, "Matthew", "Matthew"),
	asvBook(41, "Mark", "Mark"),
	asvBook(42, "Luke", "Luke"),
	asvBook(43, "John", "John"),
	asvBook(44, "Acts", "Acts"),
	asvBook(45, "Romans", "Romans"),
	asvBook(46, "1 Corinthians", "I_Corinthians"),
	asvBook(47, "2 Corinthians", "II_Corinthians"),
	asvBook(48, "Galatians", "Galatians"),
	asvBook(49, "Ephesians", "Ephesians"),
	asvBook(50, "Philippians", "Philippians"),
	asvBook(51, "Colossians", "Colossians"),
	asvBook(52, "1 Thessalonians", "I_Thessalonians"),
	asvBook(53, "2 Thessalonians", "II_Thessalonians"),
	asvBook(54, "1 Timothy", "I_Timothy"),
	asvBook(55, "2 Timothy", "II_Timothy"),
	asvBook(56, "Titus", "Titus"),
	asvBook(57, "Philemon", "Philemon"),
	asvBook(58, "Hebrews", "Hebrews"),
	asvBook(59, "James", "James"),
	asvBook(60, "1 Peter", "I_Peter"),
	asvBook(61, "2 Peter", "II_Peter"),
	asvBook(62, "1 John", "I_John"),
	asvBook(63, "2 John", "II_John"),
	asvBook(64, "3 John", "III_John"),
	asvBook(65, "Jude", "Jude"),
	asvBook(66, "Revelation", "Revelation"),
}

// ---------------------------------------------------------------------------
// Canonical Books — Greek (LXX Septuagint OT + Greek NT via Polyglot Bible)
// ---------------------------------------------------------------------------

var greekCanonicalBooks = []metadata.BookMeta{
	// ---- Old Testament (Septuagint / LXX) ----
	greekBook(1, "Genesis", "gen"),
	greekBook(2, "Exodus", "exo"),
	greekBook(3, "Leviticus", "lev"),
	greekBook(4, "Numbers", "num"),
	greekBook(5, "Deuteronomy", "deu"),
	greekBook(6, "Joshua", "jos"),
	greekBook(7, "Judges", "jdg"),
	greekBook(8, "Ruth", "rut"),
	greekBook(9, "1 Samuel", "1sa"),
	greekBook(10, "2 Samuel", "2sa"),
	greekBook(11, "1 Kings", "1ki"),
	greekBook(12, "2 Kings", "2ki"),
	greekBook(13, "1 Chronicles", "1ch"),
	greekBook(14, "2 Chronicles", "2ch"),
	greekBook(15, "Ezra", "ezr"),
	greekBook(16, "Nehemiah", "neh"),
	greekBook(17, "Esther", "est"),
	greekBook(18, "Job", "job"),
	greekBook(19, "Psalms", "psa"),
	greekBook(20, "Proverbs", "pro"),
	greekBook(21, "Ecclesiastes", "ecc"),
	greekBook(22, "Song of Solomon", "sng"),
	greekBook(23, "Isaiah", "isa"),
	greekBook(24, "Jeremiah", "jer"),
	greekBook(25, "Lamentations", "lam"),
	greekBook(26, "Ezekiel", "ezk"),
	greekBook(27, "Daniel", "dan"),
	greekBook(28, "Hosea", "hos"),
	greekBook(29, "Joel", "joe"),
	greekBook(30, "Amos", "amo"),
	greekBook(31, "Obadiah", "oba"),
	greekBook(32, "Jonah", "jon"),
	greekBook(33, "Micah", "mic"),
	greekBook(34, "Nahum", "nah"),
	greekBook(35, "Habakkuk", "hab"),
	greekBook(36, "Zephaniah", "zep"),
	greekBook(37, "Haggai", "hag"),
	greekBook(38, "Zechariah", "zec"),
	greekBook(39, "Malachi", "mal"),
	// ---- New Testament (Greek NT, Textus Receptus) ----
	greekBook(40, "Matthew", "mat"),
	greekBook(41, "Mark", "mrk"),
	greekBook(42, "Luke", "luk"),
	greekBook(43, "John", "jhn"),
	greekBook(44, "Acts", "act"),
	greekBook(45, "Romans", "rom"),
	greekBook(46, "1 Corinthians", "1co"),
	greekBook(47, "2 Corinthians", "2co"),
	greekBook(48, "Galatians", "gal"),
	greekBook(49, "Ephesians", "eph"),
	greekBook(50, "Philippians", "php"),
	greekBook(51, "Colossians", "col"),
	greekBook(52, "1 Thessalonians", "1th"),
	greekBook(53, "2 Thessalonians", "2th"),
	greekBook(54, "1 Timothy", "1ti"),
	greekBook(55, "2 Timothy", "2ti"),
	greekBook(56, "Titus", "tit"),
	greekBook(57, "Philemon", "phm"),
	greekBook(58, "Hebrews", "heb"),
	greekBook(59, "James", "jas"),
	greekBook(60, "1 Peter", "1pe"),
	greekBook(61, "2 Peter", "2pe"),
	greekBook(62, "1 John", "1jo"),
	greekBook(63, "2 John", "2jo"),
	greekBook(64, "3 John", "3jo"),
	greekBook(65, "Jude", "jud"),
	greekBook(66, "Revelation", "rev"),
}

// ---------------------------------------------------------------------------
// Canonical Books — Hebrew Masoretic Text (OT only, 39 books)
//
// Mechon Mamre book codes follow the Tanakh order:
//   01–05  Torah  |  06–11  Former Prophets  |  12–26  Latter Prophets
//   27–28  Chronicles  |  29–39  Writings
// ---------------------------------------------------------------------------

var hebrewCanonicalBooks = []metadata.BookMeta{
	hebrewBook(1, "Genesis", "01"),
	hebrewBook(2, "Exodus", "02"),
	hebrewBook(3, "Leviticus", "03"),
	hebrewBook(4, "Numbers", "04"),
	hebrewBook(5, "Deuteronomy", "05"),
	hebrewBook(6, "Joshua", "06"),
	hebrewBook(7, "Judges", "07"),
	hebrewBook(8, "Ruth", "33"),
	hebrewBook(9, "1 Samuel", "08"),
	hebrewBook(10, "2 Samuel", "09"),
	hebrewBook(11, "1 Kings", "10"),
	hebrewBook(12, "2 Kings", "11"),
	hebrewBook(13, "1 Chronicles", "27"),
	hebrewBook(14, "2 Chronicles", "28"),
	hebrewBook(15, "Ezra", "38"),
	hebrewBook(16, "Nehemiah", "39"),
	hebrewBook(17, "Esther", "36"),
	hebrewBook(18, "Job", "31"),
	hebrewBook(19, "Psalms", "29"),
	hebrewBook(20, "Proverbs", "30"),
	hebrewBook(21, "Ecclesiastes", "35"),
	hebrewBook(22, "Song of Solomon", "32"),
	hebrewBook(23, "Isaiah", "12"),
	hebrewBook(24, "Jeremiah", "13"),
	hebrewBook(25, "Lamentations", "34"),
	hebrewBook(26, "Ezekiel", "14"),
	hebrewBook(27, "Daniel", "37"),
	hebrewBook(28, "Hosea", "15"),
	hebrewBook(29, "Joel", "16"),
	hebrewBook(30, "Amos", "17"),
	hebrewBook(31, "Obadiah", "18"),
	hebrewBook(32, "Jonah", "19"),
	hebrewBook(33, "Micah", "20"),
	hebrewBook(34, "Nahum", "21"),
	hebrewBook(35, "Habakkuk", "22"),
	hebrewBook(36, "Zephaniah", "23"),
	hebrewBook(37, "Haggai", "24"),
	hebrewBook(38, "Zechariah", "25"),
	hebrewBook(39, "Malachi", "26"),
}

// ---------------------------------------------------------------------------
// Canonical Books — WLC (Westminster Leningrad Codex, same 39 books)
// ---------------------------------------------------------------------------

var wlcCanonicalBooks = []metadata.BookMeta{
	wlcBook(1, "Genesis", "01"),
	wlcBook(2, "Exodus", "02"),
	wlcBook(3, "Leviticus", "03"),
	wlcBook(4, "Numbers", "04"),
	wlcBook(5, "Deuteronomy", "05"),
	wlcBook(6, "Joshua", "06"),
	wlcBook(7, "Judges", "07"),
	wlcBook(8, "Ruth", "33"),
	wlcBook(9, "1 Samuel", "08"),
	wlcBook(10, "2 Samuel", "09"),
	wlcBook(11, "1 Kings", "10"),
	wlcBook(12, "2 Kings", "11"),
	wlcBook(13, "1 Chronicles", "27"),
	wlcBook(14, "2 Chronicles", "28"),
	wlcBook(15, "Ezra", "38"),
	wlcBook(16, "Nehemiah", "39"),
	wlcBook(17, "Esther", "36"),
	wlcBook(18, "Job", "31"),
	wlcBook(19, "Psalms", "29"),
	wlcBook(20, "Proverbs", "30"),
	wlcBook(21, "Ecclesiastes", "35"),
	wlcBook(22, "Song of Solomon", "32"),
	wlcBook(23, "Isaiah", "12"),
	wlcBook(24, "Jeremiah", "13"),
	wlcBook(25, "Lamentations", "34"),
	wlcBook(26, "Ezekiel", "14"),
	wlcBook(27, "Daniel", "37"),
	wlcBook(28, "Hosea", "15"),
	wlcBook(29, "Joel", "16"),
	wlcBook(30, "Amos", "17"),
	wlcBook(31, "Obadiah", "18"),
	wlcBook(32, "Jonah", "19"),
	wlcBook(33, "Micah", "20"),
	wlcBook(34, "Nahum", "21"),
	wlcBook(35, "Habakkuk", "22"),
	wlcBook(36, "Zephaniah", "23"),
	wlcBook(37, "Haggai", "24"),
	wlcBook(38, "Zechariah", "25"),
	wlcBook(39, "Malachi", "26"),
}

// ---------------------------------------------------------------------------
// Canonical Books — Shona (Bhaibheri Dzvene, full 66-book canon)
// ---------------------------------------------------------------------------

var shonaCanonicalBooks = []metadata.BookMeta{
	// ---- Old Testament ----
	shonaBook(1, "Genesis", "GEN"),
	shonaBook(2, "Exodus", "EXO"),
	shonaBook(3, "Leviticus", "LEV"),
	shonaBook(4, "Numbers", "NUM"),
	shonaBook(5, "Deuteronomy", "DEU"),
	shonaBook(6, "Joshua", "JOS"),
	shonaBook(7, "Judges", "JDG"),
	shonaBook(8, "Ruth", "RUT"),
	shonaBook(9, "1 Samuel", "1SA"),
	shonaBook(10, "2 Samuel", "2SA"),
	shonaBook(11, "1 Kings", "1KI"),
	shonaBook(12, "2 Kings", "2KI"),
	shonaBook(13, "1 Chronicles", "1CH"),
	shonaBook(14, "2 Chronicles", "2CH"),
	shonaBook(15, "Ezra", "EZR"),
	shonaBook(16, "Nehemiah", "NEH"),
	shonaBook(17, "Esther", "EST"),
	shonaBook(18, "Job", "JOB"),
	shonaBook(19, "Psalms", "PSA"),
	shonaBook(20, "Proverbs", "PRO"),
	shonaBook(21, "Ecclesiastes", "ECC"),
	shonaBook(22, "Song of Solomon", "SNG"),
	shonaBook(23, "Isaiah", "ISA"),
	shonaBook(24, "Jeremiah", "JER"),
	shonaBook(25, "Lamentations", "LAM"),
	shonaBook(26, "Ezekiel", "EZK"),
	shonaBook(27, "Daniel", "DAN"),
	shonaBook(28, "Hosea", "HOS"),
	shonaBook(29, "Joel", "JOL"),
	shonaBook(30, "Amos", "AMO"),
	shonaBook(31, "Obadiah", "OBA"),
	shonaBook(32, "Jonah", "JON"),
	shonaBook(33, "Micah", "MIC"),
	shonaBook(34, "Nahum", "NAH"),
	shonaBook(35, "Habakkuk", "HAB"),
	shonaBook(36, "Zephaniah", "ZEP"),
	shonaBook(37, "Haggai", "HAG"),
	shonaBook(38, "Zechariah", "ZEC"),
	shonaBook(39, "Malachi", "MAL"),
	// ---- New Testament ----
	shonaBook(40, "Matthew", "MAT"),
	shonaBook(41, "Mark", "MRK"),
	shonaBook(42, "Luke", "LUK"),
	shonaBook(43, "John", "JHN"),
	shonaBook(44, "Acts", "ACT"),
	shonaBook(45, "Romans", "ROM"),
	shonaBook(46, "1 Corinthians", "1CO"),
	shonaBook(47, "2 Corinthians", "2CO"),
	shonaBook(48, "Galatians", "GAL"),
	shonaBook(49, "Ephesians", "EPH"),
	shonaBook(50, "Philippians", "PHP"),
	shonaBook(51, "Colossians", "COL"),
	shonaBook(52, "1 Thessalonians", "1TH"),
	shonaBook(53, "2 Thessalonians", "2TH"),
	shonaBook(54, "1 Timothy", "1TI"),
	shonaBook(55, "2 Timothy", "2TI"),
	shonaBook(56, "Titus", "TIT"),
	shonaBook(57, "Philemon", "PHM"),
	shonaBook(58, "Hebrews", "HEB"),
	shonaBook(59, "James", "JAS"),
	shonaBook(60, "1 Peter", "1PE"),
	shonaBook(61, "2 Peter", "2PE"),
	shonaBook(62, "1 John", "1JO"),
	shonaBook(63, "2 John", "2JO"),
	shonaBook(64, "3 John", "3JO"),
	shonaBook(65, "Jude", "JUD"),
	shonaBook(66, "Revelation", "REV"),
}

// ---------------------------------------------------------------------------
// Canonical Books — SBL Greek New Testament (NT only, 27 books)
// ---------------------------------------------------------------------------

var sblgCanonicalBooks = []metadata.BookMeta{
	sblgBook(40, "Matthew", "Matt.txt"),
	sblgBook(41, "Mark", "Mark.txt"),
	sblgBook(42, "Luke", "Luke.txt"),
	sblgBook(43, "John", "John.txt"),
	sblgBook(44, "Acts", "Acts.txt"),
	sblgBook(45, "Romans", "Rom.txt"),
	sblgBook(46, "1 Corinthians", "1Cor.txt"),
	sblgBook(47, "2 Corinthians", "2Cor.txt"),
	sblgBook(48, "Galatians", "Gal.txt"),
	sblgBook(49, "Ephesians", "Eph.txt"),
	sblgBook(50, "Philippians", "Phil.txt"),
	sblgBook(51, "Colossians", "Col.txt"),
	sblgBook(52, "1 Thessalonians", "1Thess.txt"),
	sblgBook(53, "2 Thessalonians", "2Thess.txt"),
	sblgBook(54, "1 Timothy", "1Tim.txt"),
	sblgBook(55, "2 Timothy", "2Tim.txt"),
	sblgBook(56, "Titus", "Titus.txt"),
	sblgBook(57, "Philemon", "Phlm.txt"),
	sblgBook(58, "Hebrews", "Heb.txt"),
	sblgBook(59, "James", "Jas.txt"),
	sblgBook(60, "1 Peter", "1Pet.txt"),
	sblgBook(61, "2 Peter", "2Pet.txt"),
	sblgBook(62, "1 John", "1John.txt"),
	sblgBook(63, "2 John", "2John.txt"),
	sblgBook(64, "3 John", "3John.txt"),
	sblgBook(65, "Jude", "Jude.txt"),
	sblgBook(66, "Revelation", "Rev.txt"),
}

// ---------------------------------------------------------------------------
// Canonical Books — Young's Literal Translation (Full 66-book canon)
// ---------------------------------------------------------------------------

var yltCanonicalBooks = []metadata.BookMeta{
	yltBook(1, "Genesis"), yltBook(2, "Exodus"), yltBook(3, "Leviticus"), yltBook(4, "Numbers"), yltBook(5, "Deuteronomy"),
	yltBook(6, "Joshua"), yltBook(7, "Judges"), yltBook(8, "Ruth"), yltBook(9, "1 Samuel"), yltBook(10, "2 Samuel"),
	yltBook(11, "1 Kings"), yltBook(12, "2 Kings"), yltBook(13, "1 Chronicles"), yltBook(14, "2 Chronicles"), yltBook(15, "Ezra"),
	yltBook(16, "Nehemiah"), yltBook(17, "Esther"), yltBook(18, "Job"), yltBook(19, "Psalms"), yltBook(20, "Proverbs"),
	yltBook(21, "Ecclesiastes"), yltBook(22, "Song of Solomon"), yltBook(23, "Isaiah"), yltBook(24, "Jeremiah"), yltBook(25, "Lamentations"),
	yltBook(26, "Ezekiel"), yltBook(27, "Daniel"), yltBook(28, "Hosea"), yltBook(29, "Joel"), yltBook(30, "Amos"),
	yltBook(31, "Obadiah"), yltBook(32, "Jonah"), yltBook(33, "Micah"), yltBook(34, "Nahum"), yltBook(35, "Habakkuk"),
	yltBook(36, "Zephaniah"), yltBook(37, "Haggai"), yltBook(38, "Zechariah"), yltBook(39, "Malachi"),
	yltBook(40, "Matthew"), yltBook(41, "Mark"), yltBook(42, "Luke"), yltBook(43, "John"), yltBook(44, "Acts"),
	yltBook(45, "Romans"), yltBook(46, "1 Corinthians"), yltBook(47, "2 Corinthians"), yltBook(48, "Galatians"), yltBook(49, "Ephesians"),
	yltBook(50, "Philippians"), yltBook(51, "Colossians"), yltBook(52, "1 Thessalonians"), yltBook(53, "2 Thessalonians"), yltBook(54, "1 Timothy"),
	yltBook(55, "2 Timothy"), yltBook(56, "Titus"), yltBook(57, "Philemon"), yltBook(58, "Hebrews"), yltBook(59, "James"),
	yltBook(60, "1 Peter"), yltBook(61, "2 Peter"), yltBook(62, "1 John"), yltBook(63, "2 John"), yltBook(64, "3 John"),
	yltBook(65, "Jude"), yltBook(66, "Revelation"),
}

// ---------------------------------------------------------------------------
// Canonical Books — Darby Bible (Full 66-book canon)
// ---------------------------------------------------------------------------

var dbtCanonicalBooks = []metadata.BookMeta{
	dbtBook(1, "Genesis"), dbtBook(2, "Exodus"), dbtBook(3, "Leviticus"), dbtBook(4, "Numbers"), dbtBook(5, "Deuteronomy"),
	dbtBook(6, "Joshua"), dbtBook(7, "Judges"), dbtBook(8, "Ruth"), dbtBook(9, "1 Samuel"), dbtBook(10, "2 Samuel"),
	dbtBook(11, "1 Kings"), dbtBook(12, "2 Kings"), dbtBook(13, "1 Chronicles"), dbtBook(14, "2 Chronicles"), dbtBook(15, "Ezra"),
	dbtBook(16, "Nehemiah"), dbtBook(17, "Esther"), dbtBook(18, "Job"), dbtBook(19, "Psalms"), dbtBook(20, "Proverbs"),
	dbtBook(21, "Ecclesiastes"), dbtBook(22, "Song of Solomon"), dbtBook(23, "Isaiah"), dbtBook(24, "Jeremiah"), dbtBook(25, "Lamentations"),
	dbtBook(26, "Ezekiel"), dbtBook(27, "Daniel"), dbtBook(28, "Hosea"), dbtBook(29, "Joel"), dbtBook(30, "Amos"),
	dbtBook(31, "Obadiah"), dbtBook(32, "Jonah"), dbtBook(33, "Micah"), dbtBook(34, "Nahum"), dbtBook(35, "Habakkuk"),
	dbtBook(36, "Zephaniah"), dbtBook(37, "Haggai"), dbtBook(38, "Zechariah"), dbtBook(39, "Malachi"),
	dbtBook(40, "Matthew"), dbtBook(41, "Mark"), dbtBook(42, "Luke"), dbtBook(43, "John"), dbtBook(44, "Acts"),
	dbtBook(45, "Romans"), dbtBook(46, "1 Corinthians"), dbtBook(47, "2 Corinthians"), dbtBook(48, "Galatians"), dbtBook(49, "Ephesians"),
	dbtBook(50, "Philippians"), dbtBook(51, "Colossians"), dbtBook(52, "1 Thessalonians"), dbtBook(53, "2 Thessalonians"), dbtBook(54, "1 Timothy"),
	dbtBook(55, "2 Timothy"), dbtBook(56, "Titus"), dbtBook(57, "Philemon"), dbtBook(58, "Hebrews"), dbtBook(59, "James"),
	dbtBook(60, "1 Peter"), dbtBook(61, "2 Peter"), dbtBook(62, "1 John"), dbtBook(63, "2 John"), dbtBook(64, "3 John"),
	dbtBook(65, "Jude"), dbtBook(66, "Revelation"),
}

// ---------------------------------------------------------------------------
// Canonical Books — Douay-Rheims Bible (Full 66-book canon)
// ---------------------------------------------------------------------------

var drbCanonicalBooks = []metadata.BookMeta{
	drbBook(1, "Genesis"), drbBook(2, "Exodus"), drbBook(3, "Leviticus"), drbBook(4, "Numbers"), drbBook(5, "Deuteronomy"),
	drbBook(6, "Joshua"), drbBook(7, "Judges"), drbBook(8, "Ruth"), drbBook(9, "1 Samuel"), drbBook(10, "2 Samuel"),
	drbBook(11, "1 Kings"), drbBook(12, "2 Kings"), drbBook(13, "1 Chronicles"), drbBook(14, "2 Chronicles"), drbBook(15, "Ezra"),
	drbBook(16, "Nehemiah"), drbBook(17, "Esther"), drbBook(18, "Job"), drbBook(19, "Psalms"), drbBook(20, "Proverbs"),
	drbBook(21, "Ecclesiastes"), drbBook(22, "Song of Solomon"), drbBook(23, "Isaiah"), drbBook(24, "Jeremiah"), drbBook(25, "Lamentations"),
	drbBook(26, "Ezekiel"), drbBook(27, "Daniel"), drbBook(28, "Hosea"), drbBook(29, "Joel"), drbBook(30, "Amos"),
	drbBook(31, "Obadiah"), drbBook(32, "Jonah"), drbBook(33, "Micah"), drbBook(34, "Nahum"), drbBook(35, "Habakkuk"),
	drbBook(36, "Zephaniah"), drbBook(37, "Haggai"), drbBook(38, "Zechariah"), drbBook(39, "Malachi"),
	drbBook(40, "Matthew"), drbBook(41, "Mark"), drbBook(42, "Luke"), drbBook(43, "John"), drbBook(44, "Acts"),
	drbBook(45, "Romans"), drbBook(46, "1 Corinthians"), drbBook(47, "2 Corinthians"), drbBook(48, "Galatians"), drbBook(49, "Ephesians"),
	drbBook(50, "Philippians"), drbBook(51, "Colossians"), drbBook(52, "1 Thessalonians"), drbBook(53, "2 Thessalonians"), drbBook(54, "1 Timothy"),
	drbBook(55, "2 Timothy"), drbBook(56, "Titus"), drbBook(57, "Philemon"), drbBook(58, "Hebrews"), drbBook(59, "James"),
	drbBook(60, "1 Peter"), drbBook(61, "2 Peter"), drbBook(62, "1 John"), drbBook(63, "2 John"), drbBook(64, "3 John"),
	drbBook(65, "Jude"), drbBook(66, "Revelation"),
}

// ---------------------------------------------------------------------------
// Canonical Books — Textus Receptus Greek (NT only, 27 books)
// ---------------------------------------------------------------------------

var trCanonicalBooks = []metadata.BookMeta{
	trBook(40, "Matthew"), trBook(41, "Mark"), trBook(42, "Luke"), trBook(43, "John"), trBook(44, "Acts"),
	trBook(45, "Romans"), trBook(46, "1 Corinthians"), trBook(47, "2 Corinthians"), trBook(48, "Galatians"), trBook(49, "Ephesians"),
	trBook(50, "Philippians"), trBook(51, "Colossians"), trBook(52, "1 Thessalonians"), trBook(53, "2 Thessalonians"), trBook(54, "1 Timothy"),
	trBook(55, "2 Timothy"), trBook(56, "Titus"), trBook(57, "Philemon"), trBook(58, "Hebrews"), trBook(59, "James"),
	trBook(60, "1 Peter"), trBook(61, "2 Peter"), trBook(62, "1 John"), trBook(63, "2 John"), trBook(64, "3 John"),
	trBook(65, "Jude"), trBook(66, "Revelation"),
}

// ---------------------------------------------------------------------------
// Extra-Canonical / Deuterocanonical / Ethiopian Books (25 books)
// ---------------------------------------------------------------------------

var extraCanonicalBooks = []metadata.BookMeta{
	{
		Title:          "Enoch",
		CanonicalOrder: 1,
		Category:       metadata.CategoryExtraCanonical,
		Language:       "en",
		TranslationID:  "kjv",
		License:        "Public Domain",
		Notes:          "1 Enoch; canonical in the Ethiopian Orthodox Tewahedo Church.",
		Sources: []metadata.SourceRef{
			sacredTexts("chr/boo/enoch.htm"),
			gutenberg("30379/30379.txt"),
		},
	},
	{
		Title:          "Jubilees",
		CanonicalOrder: 2,
		Category:       metadata.CategoryExtraCanonical,
		Language:       "en",
		TranslationID:  "kjv",
		License:        "Public Domain",
		Notes:          "Book of Jubilees; canonical in the Ethiopian Orthodox Tewahedo Church.",
		Sources: []metadata.SourceRef{
			sacredTexts("chr/boo/jubilees.htm"),
		},
	},
	{
		Title:          "Tobit",
		CanonicalOrder: 3,
		Category:       metadata.CategoryExtraCanonical,
		Language:       "en",
		TranslationID:  "kjv",
		License:        "Public Domain",
		Notes:          "Deuterocanonical; accepted by Catholic and Orthodox canons.",
		Sources: []metadata.SourceRef{
			sacredTexts("chr/apo/tobit.htm"),
			wikisource("Tobit_(Douay-Rheims)"),
		},
	},
	{
		Title:          "Judith",
		CanonicalOrder: 4,
		Category:       metadata.CategoryExtraCanonical,
		Language:       "en",
		TranslationID:  "kjv",
		License:        "Public Domain",
		Notes:          "Deuterocanonical; accepted by Catholic and Orthodox canons.",
		Sources: []metadata.SourceRef{
			sacredTexts("chr/apo/judith.htm"),
			wikisource("Judith_(Douay-Rheims)"),
		},
	},
	{
		Title:          "Wisdom of Solomon",
		CanonicalOrder: 5,
		Category:       metadata.CategoryExtraCanonical,
		Language:       "en",
		TranslationID:  "kjv",
		License:        "Public Domain",
		Notes:          "Deuterocanonical; accepted by Catholic and Orthodox canons.",
		Sources: []metadata.SourceRef{
			sacredTexts("chr/apo/wisd.htm"),
			wikisource("Wisdom_of_Solomon_(Douay-Rheims)"),
		},
	},
	{
		Title:          "Sirach",
		CanonicalOrder: 6,
		Category:       metadata.CategoryExtraCanonical,
		Language:       "en",
		TranslationID:  "kjv",
		License:        "Public Domain",
		Notes:          "Ecclesiasticus; deuterocanonical book accepted by Catholic and Orthodox canons.",
		Sources: []metadata.SourceRef{
			sacredTexts("chr/apo/ecclus.htm"),
			wikisource("Ecclesiasticus_(Douay-Rheims)"),
		},
	},
	{
		Title:          "Baruch",
		CanonicalOrder: 7,
		Category:       metadata.CategoryExtraCanonical,
		Language:       "en",
		TranslationID:  "kjv",
		License:        "Public Domain",
		Notes:          "Deuterocanonical; accepted by Catholic and Orthodox canons.",
		Sources: []metadata.SourceRef{
			sacredTexts("chr/apo/baruch.htm"),
			wikisource("Baruch_(Douay-Rheims)"),
		},
	},
	{
		Title:          "Letter of Jeremiah",
		CanonicalOrder: 8,
		Category:       metadata.CategoryExtraCanonical,
		Language:       "en",
		TranslationID:  "kjv",
		License:        "Public Domain",
		Notes:          "Often treated as chapter 6 of Baruch.",
		Sources: []metadata.SourceRef{
			sacredTexts("chr/apo/LetterJeremiah.htm"),
		},
	},
	{
		Title:          "1 Esdras",
		CanonicalOrder: 9,
		Category:       metadata.CategoryExtraCanonical,
		Language:       "en",
		TranslationID:  "kjv",
		License:        "Public Domain",
		Notes:          "Also known as 3 Ezra; extra-canonical apocryphal text.",
		Sources: []metadata.SourceRef{
			sacredTexts("chr/apo/1esd.htm"),
		},
	},
	{
		Title:          "2 Esdras",
		CanonicalOrder: 10,
		Category:       metadata.CategoryExtraCanonical,
		Language:       "en",
		TranslationID:  "kjv",
		License:        "Public Domain",
		Notes:          "Also known as 4 Ezra; apocalyptic text.",
		Sources: []metadata.SourceRef{
			sacredTexts("chr/apo/2esd.htm"),
		},
	},
	{
		Title:          "1 Meqabyan",
		CanonicalOrder: 11,
		Category:       metadata.CategoryExtraCanonical,
		Language:       "en",
		TranslationID:  "kjv",
		License:        "Public Domain",
		Notes:          "Ethiopian Orthodox canon; first book of Meqabyan (distinct from 1 Maccabees).",
		Sources: []metadata.SourceRef{
			sacredTexts("chr/apo/1macc.htm"),
		},
	},
	{
		Title:          "2 Meqabyan",
		CanonicalOrder: 12,
		Category:       metadata.CategoryExtraCanonical,
		Language:       "en",
		TranslationID:  "kjv",
		License:        "Public Domain",
		Notes:          "Ethiopian Orthodox canon; second book of Meqabyan.",
		Sources: []metadata.SourceRef{
			sacredTexts("chr/apo/2macc.htm"),
		},
	},
	{
		Title:          "3 Meqabyan",
		CanonicalOrder: 13,
		Category:       metadata.CategoryExtraCanonical,
		Language:       "en",
		TranslationID:  "kjv",
		License:        "Public Domain",
		Notes:          "Ethiopian Orthodox canon; third book of Meqabyan.",
		Sources: []metadata.SourceRef{
			sacredTexts("chr/apo/3macc.htm"),
		},
	},
	{
		Title:          "Tegsats",
		CanonicalOrder: 14,
		Category:       metadata.CategoryExtraCanonical,
		Language:       "en",
		TranslationID:  "kjv",
		License:        "Public Domain",
		Notes:          "Ethiopian Orthodox canon; part of the narrower canon.",
		Sources: []metadata.SourceRef{
			{
				URL:      "https://en.wikisource.org/wiki/Portal:Ethiopian_Orthodox_canon",
				Format:   "html",
				Language: "en",
				License:  "Public Domain",
				Notes:    "Portal for Ethiopian Orthodox canonical texts.",
			},
		},
	},
	{
		Title:          "Josippon",
		CanonicalOrder: 15,
		Category:       metadata.CategoryExtraCanonical,
		Language:       "en",
		TranslationID:  "kjv",
		License:        "Public Domain",
		Notes:          "Medieval Hebrew chronicle; included in the Ethiopian Orthodox broader canon.",
		Sources: []metadata.SourceRef{
			sacredTexts("jud/josephus/index.htm"),
		},
	},
	{
		Title:          "Prayer of Manasseh",
		CanonicalOrder: 16,
		Category:       metadata.CategoryExtraCanonical,
		Language:       "en",
		TranslationID:  "kjv",
		License:        "Public Domain",
		Notes:          "Short penitential prayer; canonical in some Orthodox traditions.",
		Sources: []metadata.SourceRef{
			sacredTexts("chr/apo/manasseh.htm"),
		},
	},
	{
		Title:          "Psalm 151",
		CanonicalOrder: 17,
		Category:       metadata.CategoryExtraCanonical,
		Language:       "en",
		TranslationID:  "kjv",
		License:        "Public Domain",
		Notes:          "Found in the Septuagint and Dead Sea Scrolls; canonical in Orthodox traditions.",
		Sources: []metadata.SourceRef{
			sacredTexts("chr/apo/ps151.htm"),
		},
	},
	{
		Title:          "Sirate Tsion",
		CanonicalOrder: 18,
		Category:       metadata.CategoryExtraCanonical,
		Language:       "en",
		TranslationID:  "kjv",
		License:        "Public Domain",
		Notes:          "Ethiopian Orthodox canon; \"Order of Zion\".",
		Sources: []metadata.SourceRef{
			{
				URL:      "https://en.wikisource.org/wiki/Portal:Ethiopian_Orthodox_canon",
				Format:   "html",
				Language: "en",
				License:  "Public Domain",
				Notes:    "Ethiopian Orthodox broader canon portal.",
			},
		},
	},
	{
		Title:          "Tizaz",
		CanonicalOrder: 19,
		Category:       metadata.CategoryExtraCanonical,
		Language:       "en",
		TranslationID:  "kjv",
		License:        "Public Domain",
		Notes:          "Ethiopian Orthodox canon; part of the broader Haile Selassie Bible.",
		Sources: []metadata.SourceRef{
			{
				URL:      "https://en.wikisource.org/wiki/Portal:Ethiopian_Orthodox_canon",
				Format:   "html",
				Language: "en",
				License:  "Public Domain",
				Notes:    "Ethiopian Orthodox broader canon portal.",
			},
		},
	},
	{
		Title:          "Gitsiw",
		CanonicalOrder: 20,
		Category:       metadata.CategoryExtraCanonical,
		Language:       "en",
		TranslationID:  "kjv",
		License:        "Public Domain",
		Notes:          "Ethiopian Orthodox canon.",
		Sources: []metadata.SourceRef{
			{
				URL:      "https://en.wikisource.org/wiki/Portal:Ethiopian_Orthodox_canon",
				Format:   "html",
				Language: "en",
				License:  "Public Domain",
				Notes:    "Ethiopian Orthodox broader canon portal.",
			},
		},
	},
	{
		Title:          "Abtilis",
		CanonicalOrder: 21,
		Category:       metadata.CategoryExtraCanonical,
		Language:       "en",
		TranslationID:  "kjv",
		License:        "Public Domain",
		Notes:          "Ethiopian Orthodox canon.",
		Sources: []metadata.SourceRef{
			{
				URL:      "https://en.wikisource.org/wiki/Portal:Ethiopian_Orthodox_canon",
				Format:   "html",
				Language: "en",
				License:  "Public Domain",
				Notes:    "Ethiopian Orthodox broader canon portal.",
			},
		},
	},
	{
		Title:          "Metsihafe Kidan I",
		CanonicalOrder: 22,
		Category:       metadata.CategoryExtraCanonical,
		Language:       "en",
		TranslationID:  "kjv",
		License:        "Public Domain",
		Notes:          "Ethiopian Orthodox canon; \"Book of the Covenant\" part I.",
		Sources: []metadata.SourceRef{
			{
				URL:      "https://en.wikisource.org/wiki/Portal:Ethiopian_Orthodox_canon",
				Format:   "html",
				Language: "en",
				License:  "Public Domain",
				Notes:    "Ethiopian Orthodox broader canon portal.",
			},
		},
	},
	{
		Title:          "Metsihafe Kidan II",
		CanonicalOrder: 23,
		Category:       metadata.CategoryExtraCanonical,
		Language:       "en",
		TranslationID:  "kjv",
		License:        "Public Domain",
		Notes:          "Ethiopian Orthodox canon; \"Book of the Covenant\" part II.",
		Sources: []metadata.SourceRef{
			{
				URL:      "https://en.wikisource.org/wiki/Portal:Ethiopian_Orthodox_canon",
				Format:   "html",
				Language: "en",
				License:  "Public Domain",
				Notes:    "Ethiopian Orthodox broader canon portal.",
			},
		},
	},
	{
		Title:          "Qalëmentos",
		CanonicalOrder: 24,
		Category:       metadata.CategoryExtraCanonical,
		Language:       "en",
		TranslationID:  "kjv",
		License:        "Public Domain",
		Notes:          "Ethiopian Orthodox canon; Clementine literature.",
		Sources: []metadata.SourceRef{
			{
				URL:      "https://en.wikisource.org/wiki/Portal:Ethiopian_Orthodox_canon",
				Format:   "html",
				Language: "en",
				License:  "Public Domain",
				Notes:    "Ethiopian Orthodox broader canon portal.",
			},
		},
	},
	{
		Title:          "Didesqelya",
		CanonicalOrder: 25,
		Category:       metadata.CategoryExtraCanonical,
		Language:       "en",
		TranslationID:  "kjv",
		License:        "Public Domain",
		Notes:          "Ethiopian Orthodox canon; related to the Didascalia Apostolorum.",
		Sources: []metadata.SourceRef{
			sacredTexts("chr/ecf/007/0070193.htm"),
		},
	},
}

// ---------------------------------------------------------------------------
// Dead Sea Scrolls (9 major English-translated texts)
// ---------------------------------------------------------------------------

var deadSeaScrollsBooks = []metadata.BookMeta{
	{
		Title:          "Introduction to the Dead Sea Scrolls",
		CanonicalOrder: 1,
		Category:       metadata.CategoryDeadSeaScrolls,
		Language:       "en",
		TranslationID:  "kjv",
		License:        "Public Domain",
		Notes:          "Overview and introduction to the Dead Sea Scrolls corpus.",
		Sources: []metadata.SourceRef{
			sacredTexts("jud/dss/dss01.htm"),
		},
	},
	{
		Title:          "Community Rule",
		CanonicalOrder: 2,
		Category:       metadata.CategoryDeadSeaScrolls,
		Language:       "en",
		TranslationID:  "kjv",
		License:        "Public Domain",
		Notes:          "1QS — Serekh ha-Yahad; foundational Qumran community document.",
		Sources: []metadata.SourceRef{
			sacredTexts("jud/dss/dss03.htm"),
		},
	},
	{
		Title:          "War Scroll",
		CanonicalOrder: 3,
		Category:       metadata.CategoryDeadSeaScrolls,
		Language:       "en",
		TranslationID:  "kjv",
		License:        "Public Domain",
		Notes:          "1QM — describes eschatological war between Sons of Light and Sons of Darkness.",
		Sources: []metadata.SourceRef{
			sacredTexts("jud/dss/dss04.htm"),
		},
	},
	{
		Title:          "Thanksgiving Hymns",
		CanonicalOrder: 4,
		Category:       metadata.CategoryDeadSeaScrolls,
		Language:       "en",
		TranslationID:  "kjv",
		License:        "Public Domain",
		Notes:          "1QH — Hodayot; collection of hymns from Qumran.",
		Sources: []metadata.SourceRef{
			sacredTexts("jud/dss/dss05.htm"),
		},
	},
	{
		Title:          "Temple Scroll",
		CanonicalOrder: 5,
		Category:       metadata.CategoryDeadSeaScrolls,
		Language:       "en",
		TranslationID:  "kjv",
		License:        "Public Domain",
		Notes:          "11QT — longest Dead Sea Scroll; describes an ideal Temple and laws.",
		Sources: []metadata.SourceRef{
			sacredTexts("jud/dss/dss08.htm"),
		},
	},
	{
		Title:          "Genesis Apocryphon",
		CanonicalOrder: 6,
		Category:       metadata.CategoryDeadSeaScrolls,
		Language:       "en",
		TranslationID:  "kjv",
		License:        "Public Domain",
		Notes:          "1QapGen — Aramaic retelling and expansion of Genesis narratives.",
		Sources: []metadata.SourceRef{
			sacredTexts("jud/dss/dss07.htm"),
		},
	},
	{
		Title:          "Damascus Document",
		CanonicalOrder: 7,
		Category:       metadata.CategoryDeadSeaScrolls,
		Language:       "en",
		TranslationID:  "kjv",
		License:        "Public Domain",
		Notes:          "CD — legal text governing community life; also found in Cairo Geniza.",
		Sources: []metadata.SourceRef{
			sacredTexts("jud/dss/dss06.htm"),
		},
	},
	{
		Title:          "Habakkuk Commentary",
		CanonicalOrder: 8,
		Category:       metadata.CategoryDeadSeaScrolls,
		Language:       "en",
		TranslationID:  "kjv",
		License:        "Public Domain",
		Notes:          "1QpHab — Pesher Habakkuk; commentary on the book of Habakkuk.",
		Sources: []metadata.SourceRef{
			sacredTexts("jud/dss/dss09.htm"),
		},
	},
	{
		Title:          "Some Precepts of the Law",
		CanonicalOrder: 9,
		Category:       metadata.CategoryDeadSeaScrolls,
		Language:       "en",
		TranslationID:  "kjv",
		License:        "Public Domain",
		Notes:          "4QMMT — Miqsat Ma'ase ha-Torah; halakhic letter outlining legal disagreements.",
		Sources: []metadata.SourceRef{
			sacredTexts("jud/dss/dss12.htm"),
		},
	},
}

// ---------------------------------------------------------------------------
// Public accessors — original catalogue
// ---------------------------------------------------------------------------

// CanonicalBooks returns all 66 Protestant canonical books across ALL translations.
func CanonicalBooks() []metadata.BookMeta {
	var out []metadata.BookMeta
	out = append(out, kjvCanonicalBooks...)
	out = append(out, webCanonicalBooks...)
	out = append(out, asvCanonicalBooks...)
	out = append(out, greekCanonicalBooks...)
	out = append(out, hebrewCanonicalBooks...)
	out = append(out, shonaCanonicalBooks...)
	out = append(out, wlcCanonicalBooks...)
	out = append(out, sblgCanonicalBooks...)
	out = append(out, yltCanonicalBooks...)
	out = append(out, dbtCanonicalBooks...)
	out = append(out, drbCanonicalBooks...)
	out = append(out, trCanonicalBooks...)
	return out
}

// ExtraCanonicalBooks returns all extra-canonical / deuterocanonical books.
func ExtraCanonicalBooks() []metadata.BookMeta {
	out := make([]metadata.BookMeta, len(extraCanonicalBooks))
	copy(out, extraCanonicalBooks)
	return out
}

// DeadSeaScrollBooks returns all Dead Sea Scrolls texts.
func DeadSeaScrollBooks() []metadata.BookMeta {
	out := make([]metadata.BookMeta, len(deadSeaScrollsBooks))
	copy(out, deadSeaScrollsBooks)
	return out
}

// AllBooks returns the complete combined catalogue (all categories, all translations).
func AllBooks() []metadata.BookMeta {
	var all []metadata.BookMeta
	all = append(all, CanonicalBooks()...)
	all = append(all, ExtraCanonicalBooks()...)
	all = append(all, DeadSeaScrollBooks()...)
	return all
}

// ---------------------------------------------------------------------------
// Public accessors — filtering by translation or language
// ---------------------------------------------------------------------------

// BooksByTranslation returns all books whose TranslationID matches id.
// Use this to build (or re-build) a single translation without touching others.
func BooksByTranslation(id string) []metadata.BookMeta {
	all := AllBooks()
	var out []metadata.BookMeta
	for _, b := range all {
		if b.TranslationID == id {
			out = append(out, b)
		}
	}
	return out
}

// BooksByLanguage returns all books whose Language code matches lang (BCP 47).
// Use this to build (or re-build) all books for a specific language.
func BooksByLanguage(lang string) []metadata.BookMeta {
	all := AllBooks()
	var out []metadata.BookMeta
	for _, b := range all {
		if b.Language == lang {
			out = append(out, b)
		}
	}
	return out
}

// AllTranslationIDs returns the deduplicated list of every known TranslationID.
func AllTranslationIDs() []string {
	seen := make(map[string]bool)
	var ids []string
	for _, b := range AllBooks() {
		if b.TranslationID != "" && !seen[b.TranslationID] {
			seen[b.TranslationID] = true
			ids = append(ids, b.TranslationID)
		}
	}
	return ids
}

// AllLanguageCodes returns the deduplicated list of every known language code.
func AllLanguageCodes() []string {
	seen := make(map[string]bool)
	var codes []string
	for _, b := range AllBooks() {
		if b.Language != "" && !seen[b.Language] {
			seen[b.Language] = true
			codes = append(codes, b.Language)
		}
	}
	return codes
}

// FindBookByOrder returns the metadata for a canonical book by its order (1-66).
// It returns only the first match (usually KJV) to get an English title.
func FindBookByOrder(order int) (metadata.BookMeta, bool) {
	for _, b := range CanonicalBooks() {
		if b.CanonicalOrder == order {
			return b, true
		}
	}
	return metadata.BookMeta{}, false
}

