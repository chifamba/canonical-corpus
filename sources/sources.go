// Package sources defines the full catalogue of biblical corpus books and their
// source URLs, covering English ESV and NIV translations plus the nearest
// scholarly equivalent in each supported language.
//
// English translations:
//
//	ESV – English Standard Version (api.esv.org, requires ESV_API_KEY)
//	NIV – New International Version (biblegateway.com / scripture.api.bible, requires BIBLE_API_KEY)
//
// Nearest equivalents in other languages:
//
//	Greek (el)  – Septuagint (LXX) for OT; SBLGNT polyglot for NT  (sacred-texts.com, public domain)
//	Hebrew (he) – Westminster Leningrad Codex (WLC) for OT           (tanach.us, public domain)
//	Shona (sn)  – Union Shona Bible (BSZ / Bhaibheri Dzvene)         (biblegateway.com)
//	Amharic (am)– Ethiopian Amharic Bible (AB)                       (biblegateway.com)
//	              Nearest modern equivalent to the Ge'ez liturgical tradition
package sources

import (
	"net/url"

	"github.com/chifamba/canonical-corpus/internal/metadata"
)

// ---------------------------------------------------------------------------
// Source-reference constructors
// ---------------------------------------------------------------------------

// esvBook returns an ESV API source reference for the named book.
// Authentication uses the ESV_API_KEY environment variable.
func esvBook(bookName string) metadata.SourceRef {
	u := "https://api.esv.org/v3/passage/text/?q=" + url.QueryEscape(bookName) +
		"&include-headings=false&include-footnotes=false" +
		"&include-short-copyright=false&include-passage-references=false"
	return metadata.SourceRef{
		URL:         u,
		Format:      "json",
		Language:    "en",
		Translation: "ESV",
		License:     "ESV® Bible (The Holy Bible, English Standard Version®) Copyright © 2001 by Crossway",
		AuthHeader:  "Authorization",
		AuthEnvVar:  "ESV_API_KEY",
		AuthPrefix:  "Token ",
	}
}

// nivBook returns a NIV source reference via Scripture API Bible.
// Authentication uses the BIBLE_API_KEY environment variable.
// NIV Bible ID on scripture.api.bible: 78a9f6124f344018-01
func nivBook(bookID string) metadata.SourceRef {
	u := "https://api.scripture.api.bible/v1/bibles/78a9f6124f344018-01/books/" +
		url.PathEscape(bookID) + "?include-chapters=true&content-type=text&include-notes=false"
	return metadata.SourceRef{
		URL:         u,
		Format:      "json",
		Language:    "en",
		Translation: "NIV",
		License:     "Holy Bible, New International Version®, NIV® Copyright © 1973–2011 by Biblica, Inc.®",
		AuthHeader:  "api-key",
		AuthEnvVar:  "BIBLE_API_KEY",
	}
}

// lxxBook returns a Septuagint (Greek OT) source reference from sacred-texts.com.
func lxxBook(stPath string) metadata.SourceRef {
	return metadata.SourceRef{
		URL:         "https://www.sacred-texts.com/bib/sep/" + stPath,
		Format:      "html",
		Language:    "el",
		Translation: "LXX",
		License:     "Public Domain",
		Notes:       "Septuagint (Greek Old Testament); Brenton's translation/Greek text.",
	}
}

// sblgntBook returns an SBLGNT/Greek-NT polyglot source reference from sacred-texts.com.
func sblgntBook(stPath string) metadata.SourceRef {
	return metadata.SourceRef{
		URL:         "https://www.sacred-texts.com/bib/poly/" + stPath,
		Format:      "html",
		Language:    "el",
		Translation: "SBLGNT",
		License:     "Public Domain",
		Notes:       "Greek New Testament (polyglot page); nearest equivalent to ESV/NIV critical text.",
	}
}

// wlcBook returns a Westminster Leningrad Codex (Hebrew OT) source from tanach.us.
func wlcBook(tanachBook string) metadata.SourceRef {
	return metadata.SourceRef{
		URL:         "https://www.tanach.us/Books/" + url.PathEscape(tanachBook) + ".xml",
		Format:      "xml",
		Language:    "he",
		Translation: "WLC",
		License:     "Public Domain",
		Notes:       "Westminster Leningrad Codex; the standard scholarly Hebrew Old Testament.",
	}
}

// shonaBook returns a Union Shona Bible (BSZ / Bhaibheri Dzvene) source from BibleGateway.
func shonaBook(bgName string) metadata.SourceRef {
	return metadata.SourceRef{
		URL:         "https://www.biblegateway.com/passage/?search=" + url.QueryEscape(bgName) + "&version=BSZ",
		Format:      "html",
		Language:    "sn",
		Translation: "BSZ",
		License:     "Contact Bible Society of Zimbabwe for licensing.",
		Notes:       "Bhaibheri Dzvene (Union Shona Bible); most widely used Shona Bible translation.",
	}
}

// amharicBook returns an Ethiopian Amharic Bible (AB) source from BibleGateway.
// The Amharic Bible is the nearest modern equivalent to the Ge'ez liturgical tradition.
func amharicBook(bgName string) metadata.SourceRef {
	return metadata.SourceRef{
		URL:         "https://www.biblegateway.com/passage/?search=" + url.QueryEscape(bgName) + "&version=AB",
		Format:      "html",
		Language:    "am",
		Translation: "AB",
		License:     "Contact Ethiopian Bible Society for licensing.",
		Notes:       "Ethiopian Amharic Bible; nearest modern equivalent to the Ge'ez liturgical tradition.",
	}
}

// ---------------------------------------------------------------------------
// Canonical book constructors
// ---------------------------------------------------------------------------

// otBook creates all translation entries for a single Old-Testament canonical book.
// It returns one BookMeta per language/translation: ESV, NIV, Greek (LXX), Hebrew (WLC),
// Shona (BSZ), and Amharic (AB).
func otBook(
	order int, title, esvName, nivID, lxxPath, wlcName, bgName string,
) []metadata.BookMeta {
	return []metadata.BookMeta{
		{
			Title:          title,
			CanonicalOrder: order,
			Category:       metadata.CategoryCanonical,
			Language:       "en",
			Translation:    "ESV",
			License:        "See source",
			Sources:        []metadata.SourceRef{esvBook(esvName)},
		},
		{
			Title:          title,
			CanonicalOrder: order,
			Category:       metadata.CategoryCanonical,
			Language:       "en",
			Translation:    "NIV",
			License:        "See source",
			Sources:        []metadata.SourceRef{nivBook(nivID)},
		},
		{
			Title:          title,
			CanonicalOrder: order,
			Category:       metadata.CategoryCanonical,
			Language:       "el",
			Translation:    "LXX",
			License:        "Public Domain",
			Sources:        []metadata.SourceRef{lxxBook(lxxPath)},
		},
		{
			Title:          title,
			CanonicalOrder: order,
			Category:       metadata.CategoryCanonical,
			Language:       "he",
			Translation:    "WLC",
			License:        "Public Domain",
			Sources:        []metadata.SourceRef{wlcBook(wlcName)},
		},
		{
			Title:          title,
			CanonicalOrder: order,
			Category:       metadata.CategoryCanonical,
			Language:       "sn",
			Translation:    "BSZ",
			License:        "See source",
			Sources:        []metadata.SourceRef{shonaBook(bgName)},
		},
		{
			Title:          title,
			CanonicalOrder: order,
			Category:       metadata.CategoryCanonical,
			Language:       "am",
			Translation:    "AB",
			License:        "See source",
			Sources:        []metadata.SourceRef{amharicBook(bgName)},
		},
	}
}

// ntBook creates all translation entries for a single New-Testament canonical book.
// It returns one BookMeta per language/translation: ESV, NIV, Greek (SBLGNT),
// Shona (BSZ), and Amharic (AB).  No Hebrew source is included for NT books.
func ntBook(
	order int, title, esvName, nivID, sblgntPath, bgName string,
) []metadata.BookMeta {
	return []metadata.BookMeta{
		{
			Title:          title,
			CanonicalOrder: order,
			Category:       metadata.CategoryCanonical,
			Language:       "en",
			Translation:    "ESV",
			License:        "See source",
			Sources:        []metadata.SourceRef{esvBook(esvName)},
		},
		{
			Title:          title,
			CanonicalOrder: order,
			Category:       metadata.CategoryCanonical,
			Language:       "en",
			Translation:    "NIV",
			License:        "See source",
			Sources:        []metadata.SourceRef{nivBook(nivID)},
		},
		{
			Title:          title,
			CanonicalOrder: order,
			Category:       metadata.CategoryCanonical,
			Language:       "el",
			Translation:    "SBLGNT",
			License:        "Public Domain",
			Sources:        []metadata.SourceRef{sblgntBook(sblgntPath)},
		},
		{
			Title:          title,
			CanonicalOrder: order,
			Category:       metadata.CategoryCanonical,
			Language:       "sn",
			Translation:    "BSZ",
			License:        "See source",
			Sources:        []metadata.SourceRef{shonaBook(bgName)},
		},
		{
			Title:          title,
			CanonicalOrder: order,
			Category:       metadata.CategoryCanonical,
			Language:       "am",
			Translation:    "AB",
			License:        "See source",
			Sources:        []metadata.SourceRef{amharicBook(bgName)},
		},
	}
}

// ---------------------------------------------------------------------------
// Canonical Books (Protestant 66-book canon)
//
// For each book:
//   esvName   – book name as accepted by api.esv.org  (usually plain English)
//   nivID     – book ID for scripture.api.bible NIV    (OSIS/USFM abbreviation)
//   lxxPath   – path on sacred-texts.com/bib/sep/      (OT only)
//   wlcName   – book name on tanach.us/Books/           (OT only)
//   sblgntPath– path on sacred-texts.com/bib/poly/      (NT only)
//   bgName    – book name for BibleGateway               (Shona & Amharic)
// ---------------------------------------------------------------------------

// collect flattens multiple otBook/ntBook slices into a single slice.
func collect(groups ...[]metadata.BookMeta) []metadata.BookMeta {
	var out []metadata.BookMeta
	for _, g := range groups {
		out = append(out, g...)
	}
	return out
}

var canonicalBooks = collect(
	// ── Old Testament ──────────────────────────────────────────────────────
	//
	//	order  title               esvName           nivID   lxxPath      wlcName        bgName
	otBook(1, "Genesis", "Genesis", "GEN", "gen.htm", "Genesis", "Genesis"),
	otBook(2, "Exodus", "Exodus", "EXO", "exo.htm", "Exodus", "Exodus"),
	otBook(3, "Leviticus", "Leviticus", "LEV", "lev.htm", "Leviticus", "Leviticus"),
	otBook(4, "Numbers", "Numbers", "NUM", "num.htm", "Numbers", "Numbers"),
	otBook(5, "Deuteronomy", "Deuteronomy", "DEU", "deu.htm", "Deuteronomy", "Deuteronomy"),
	otBook(6, "Joshua", "Joshua", "JOS", "jos.htm", "Joshua", "Joshua"),
	otBook(7, "Judges", "Judges", "JDG", "jdg.htm", "Judges", "Judges"),
	otBook(8, "Ruth", "Ruth", "RUT", "rut.htm", "Ruth", "Ruth"),
	otBook(9, "1 Samuel", "1 Samuel", "1SA", "1sa.htm", "I_Samuel", "1 Samuel"),
	otBook(10, "2 Samuel", "2 Samuel", "2SA", "2sa.htm", "II_Samuel", "2 Samuel"),
	otBook(11, "1 Kings", "1 Kings", "1KI", "1ki.htm", "I_Kings", "1 Kings"),
	otBook(12, "2 Kings", "2 Kings", "2KI", "2ki.htm", "II_Kings", "2 Kings"),
	otBook(13, "1 Chronicles", "1 Chronicles", "1CH", "1ch.htm", "I_Chronicles", "1 Chronicles"),
	otBook(14, "2 Chronicles", "2 Chronicles", "2CH", "2ch.htm", "II_Chronicles", "2 Chronicles"),
	otBook(15, "Ezra", "Ezra", "EZR", "ezr.htm", "Ezra", "Ezra"),
	otBook(16, "Nehemiah", "Nehemiah", "NEH", "neh.htm", "Nehemiah", "Nehemiah"),
	otBook(17, "Esther", "Esther", "EST", "est.htm", "Esther", "Esther"),
	otBook(18, "Job", "Job", "JOB", "job.htm", "Job", "Job"),
	otBook(19, "Psalms", "Psalms", "PSA", "psa.htm", "Psalms", "Psalms"),
	otBook(20, "Proverbs", "Proverbs", "PRO", "pro.htm", "Proverbs", "Proverbs"),
	otBook(21, "Ecclesiastes", "Ecclesiastes", "ECC", "ecc.htm", "Ecclesiastes", "Ecclesiastes"),
	otBook(22, "Song of Songs", "Song of Solomon", "SNG", "sol.htm", "Song_of_Songs", "Song of Songs"),
	otBook(23, "Isaiah", "Isaiah", "ISA", "isa.htm", "Isaiah", "Isaiah"),
	otBook(24, "Jeremiah", "Jeremiah", "JER", "jer.htm", "Jeremiah", "Jeremiah"),
	otBook(25, "Lamentations", "Lamentations", "LAM", "lam.htm", "Lamentations", "Lamentations"),
	otBook(26, "Ezekiel", "Ezekiel", "EZK", "eze.htm", "Ezekiel", "Ezekiel"),
	otBook(27, "Daniel", "Daniel", "DAN", "dan.htm", "Daniel", "Daniel"),
	otBook(28, "Hosea", "Hosea", "HOS", "hos.htm", "Hosea", "Hosea"),
	otBook(29, "Joel", "Joel", "JOL", "joe.htm", "Joel", "Joel"),
	otBook(30, "Amos", "Amos", "AMO", "amo.htm", "Amos", "Amos"),
	otBook(31, "Obadiah", "Obadiah", "OBA", "oba.htm", "Obadiah", "Obadiah"),
	otBook(32, "Jonah", "Jonah", "JON", "jon.htm", "Jonah", "Jonah"),
	otBook(33, "Micah", "Micah", "MIC", "mic.htm", "Micah", "Micah"),
	otBook(34, "Nahum", "Nahum", "NAM", "nah.htm", "Nahum", "Nahum"),
	otBook(35, "Habakkuk", "Habakkuk", "HAB", "hab.htm", "Habakkuk", "Habakkuk"),
	otBook(36, "Zephaniah", "Zephaniah", "ZEP", "zep.htm", "Zephaniah", "Zephaniah"),
	otBook(37, "Haggai", "Haggai", "HAG", "hag.htm", "Haggai", "Haggai"),
	otBook(38, "Zechariah", "Zechariah", "ZEC", "zec.htm", "Zechariah", "Zechariah"),
	otBook(39, "Malachi", "Malachi", "MAL", "mal.htm", "Malachi", "Malachi"),

	// ── New Testament ──────────────────────────────────────────────────────
	//
	//	order  title              esvName        nivID   sblgntPath  bgName
	ntBook(40, "Matthew", "Matthew", "MAT", "mat.htm", "Matthew"),
	ntBook(41, "Mark", "Mark", "MRK", "mar.htm", "Mark"),
	ntBook(42, "Luke", "Luke", "LUK", "luk.htm", "Luke"),
	ntBook(43, "John", "John", "JHN", "joh.htm", "John"),
	ntBook(44, "Acts", "Acts", "ACT", "act.htm", "Acts"),
	ntBook(45, "Romans", "Romans", "ROM", "rom.htm", "Romans"),
	ntBook(46, "1 Corinthians", "1 Corinthians", "1CO", "1co.htm", "1 Corinthians"),
	ntBook(47, "2 Corinthians", "2 Corinthians", "2CO", "2co.htm", "2 Corinthians"),
	ntBook(48, "Galatians", "Galatians", "GAL", "gal.htm", "Galatians"),
	ntBook(49, "Ephesians", "Ephesians", "EPH", "eph.htm", "Ephesians"),
	ntBook(50, "Philippians", "Philippians", "PHP", "phi.htm", "Philippians"),
	ntBook(51, "Colossians", "Colossians", "COL", "col.htm", "Colossians"),
	ntBook(52, "1 Thessalonians", "1 Thessalonians", "1TH", "1th.htm", "1 Thessalonians"),
	ntBook(53, "2 Thessalonians", "2 Thessalonians", "2TH", "2th.htm", "2 Thessalonians"),
	ntBook(54, "1 Timothy", "1 Timothy", "1TI", "1ti.htm", "1 Timothy"),
	ntBook(55, "2 Timothy", "2 Timothy", "2TI", "2ti.htm", "2 Timothy"),
	ntBook(56, "Titus", "Titus", "TIT", "tit.htm", "Titus"),
	ntBook(57, "Philemon", "Philemon", "PHM", "phm.htm", "Philemon"),
	ntBook(58, "Hebrews", "Hebrews", "HEB", "heb.htm", "Hebrews"),
	ntBook(59, "James", "James", "JAS", "jam.htm", "James"),
	ntBook(60, "1 Peter", "1 Peter", "1PE", "1pe.htm", "1 Peter"),
	ntBook(61, "2 Peter", "2 Peter", "2PE", "2pe.htm", "2 Peter"),
	ntBook(62, "1 John", "1 John", "1JN", "1jo.htm", "1 John"),
	ntBook(63, "2 John", "2 John", "2JN", "2jo.htm", "2 John"),
	ntBook(64, "3 John", "3 John", "3JN", "3jo.htm", "3 John"),
	ntBook(65, "Jude", "Jude", "JUD", "jud.htm", "Jude"),
	ntBook(66, "Revelation", "Revelation", "REV", "rev.htm", "Revelation"),
)

// ---------------------------------------------------------------------------
// Extra-Canonical / Deuterocanonical / Ethiopian Books (25 books)
// ---------------------------------------------------------------------------

func sacredTexts(path string) metadata.SourceRef {
	return metadata.SourceRef{
		URL:      "https://www.sacred-texts.com/" + path,
		Format:   "html",
		Language: "en",
		License:  "Public Domain",
	}
}

func wikisource(path string) metadata.SourceRef {
	return metadata.SourceRef{
		URL:      "https://en.wikisource.org/wiki/" + path,
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

func ethiopianOrthodoxPortal() metadata.SourceRef {
	return metadata.SourceRef{
		URL:      "https://en.wikisource.org/wiki/Portal:Ethiopian_Orthodox_canon",
		Format:   "html",
		Language: "en",
		License:  "Public Domain",
		Notes:    "Ethiopian Orthodox broader canon portal.",
	}
}

var extraCanonicalBooks = []metadata.BookMeta{
	{
		Title:          "Enoch",
		CanonicalOrder: 1,
		Category:       metadata.CategoryExtraCanonical,
		Language:       "en",
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
		License:        "Public Domain",
		Notes:          "Ethiopian Orthodox canon; part of the narrower canon.",
		Sources:        []metadata.SourceRef{ethiopianOrthodoxPortal()},
	},
	{
		Title:          "Josippon",
		CanonicalOrder: 15,
		Category:       metadata.CategoryExtraCanonical,
		Language:       "en",
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
		License:        "Public Domain",
		Notes:          "Ethiopian Orthodox canon; \"Order of Zion\".",
		Sources:        []metadata.SourceRef{ethiopianOrthodoxPortal()},
	},
	{
		Title:          "Tizaz",
		CanonicalOrder: 19,
		Category:       metadata.CategoryExtraCanonical,
		Language:       "en",
		License:        "Public Domain",
		Notes:          "Ethiopian Orthodox canon; part of the broader Haile Selassie Bible.",
		Sources:        []metadata.SourceRef{ethiopianOrthodoxPortal()},
	},
	{
		Title:          "Gitsiw",
		CanonicalOrder: 20,
		Category:       metadata.CategoryExtraCanonical,
		Language:       "en",
		License:        "Public Domain",
		Notes:          "Ethiopian Orthodox canon.",
		Sources:        []metadata.SourceRef{ethiopianOrthodoxPortal()},
	},
	{
		Title:          "Abtilis",
		CanonicalOrder: 21,
		Category:       metadata.CategoryExtraCanonical,
		Language:       "en",
		License:        "Public Domain",
		Notes:          "Ethiopian Orthodox canon.",
		Sources:        []metadata.SourceRef{ethiopianOrthodoxPortal()},
	},
	{
		Title:          "Metsihafe Kidan I",
		CanonicalOrder: 22,
		Category:       metadata.CategoryExtraCanonical,
		Language:       "en",
		License:        "Public Domain",
		Notes:          "Ethiopian Orthodox canon; \"Book of the Covenant\" part I.",
		Sources:        []metadata.SourceRef{ethiopianOrthodoxPortal()},
	},
	{
		Title:          "Metsihafe Kidan II",
		CanonicalOrder: 23,
		Category:       metadata.CategoryExtraCanonical,
		Language:       "en",
		License:        "Public Domain",
		Notes:          "Ethiopian Orthodox canon; \"Book of the Covenant\" part II.",
		Sources:        []metadata.SourceRef{ethiopianOrthodoxPortal()},
	},
	{
		Title:          "Qalëmentos",
		CanonicalOrder: 24,
		Category:       metadata.CategoryExtraCanonical,
		Language:       "en",
		License:        "Public Domain",
		Notes:          "Ethiopian Orthodox canon; Clementine literature.",
		Sources:        []metadata.SourceRef{ethiopianOrthodoxPortal()},
	},
	{
		Title:          "Didesqelya",
		CanonicalOrder: 25,
		Category:       metadata.CategoryExtraCanonical,
		Language:       "en",
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
		License:        "Public Domain",
		Notes:          "4QMMT — Miqsat Ma'ase ha-Torah; halakhic letter outlining legal disagreements.",
		Sources: []metadata.SourceRef{
			sacredTexts("jud/dss/dss12.htm"),
		},
	},
}

// ---------------------------------------------------------------------------
// Public accessors
// ---------------------------------------------------------------------------

// CanonicalBooks returns all canonical book entries (multiple per book — one per
// language/translation: ESV, NIV, Greek, Hebrew for OT, Shona, Amharic).
func CanonicalBooks() []metadata.BookMeta {
	out := make([]metadata.BookMeta, len(canonicalBooks))
	copy(out, canonicalBooks)
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

// AllBooks returns the complete combined catalogue.
func AllBooks() []metadata.BookMeta {
	var all []metadata.BookMeta
	all = append(all, CanonicalBooks()...)
	all = append(all, ExtraCanonicalBooks()...)
	all = append(all, DeadSeaScrollBooks()...)
	return all
}
