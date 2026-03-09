// Package sources defines the full catalogue of biblical corpus books and their
// public-domain source URLs.
package sources

import "github.com/chifamba/canonical-corpus/internal/metadata"

// helper constructors -----------------------------------------------------------

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

// canonical returns a canonical KJV source via Wikisource.
func canonical(order int, title, wsPath string, extra ...metadata.SourceRef) metadata.BookMeta {
	srcs := []metadata.SourceRef{wikisource(wsPath)}
	srcs = append(srcs, extra...)
	return metadata.BookMeta{
		Title:          title,
		CanonicalOrder: order,
		Category:       metadata.CategoryCanonical,
		Language:       "en",
		License:        "Public Domain",
		Sources:        srcs,
	}
}

// ---------------------------------------------------------------------------
// Canonical Books (Protestant 66-book canon, KJV order)
// ---------------------------------------------------------------------------

var canonicalBooks = []metadata.BookMeta{
	// ---- Old Testament ----
	canonical(1, "Genesis", "Bible_(King_James)/Genesis"),
	canonical(2, "Exodus", "Bible_(King_James)/Exodus"),
	canonical(3, "Leviticus", "Bible_(King_James)/Leviticus"),
	canonical(4, "Numbers", "Bible_(King_James)/Numbers"),
	canonical(5, "Deuteronomy", "Bible_(King_James)/Deuteronomy"),
	canonical(6, "Joshua", "Bible_(King_James)/Joshua"),
	canonical(7, "Judges", "Bible_(King_James)/Judges"),
	canonical(8, "Ruth", "Bible_(King_James)/Ruth"),
	canonical(9, "1 Samuel", "Bible_(King_James)/I_Samuel"),
	canonical(10, "2 Samuel", "Bible_(King_James)/II_Samuel"),
	canonical(11, "1 Kings", "Bible_(King_James)/I_Kings"),
	canonical(12, "2 Kings", "Bible_(King_James)/II_Kings"),
	canonical(13, "1 Chronicles", "Bible_(King_James)/I_Chronicles"),
	canonical(14, "2 Chronicles", "Bible_(King_James)/II_Chronicles"),
	canonical(15, "Ezra", "Bible_(King_James)/Ezra"),
	canonical(16, "Nehemiah", "Bible_(King_James)/Nehemiah"),
	canonical(17, "Esther", "Bible_(King_James)/Esther"),
	canonical(18, "Job", "Bible_(King_James)/Job"),
	canonical(19, "Psalms", "Bible_(King_James)/Psalms"),
	canonical(20, "Proverbs", "Bible_(King_James)/Proverbs"),
	canonical(21, "Ecclesiastes", "Bible_(King_James)/Ecclesiastes"),
	canonical(22, "Song of Solomon", "Bible_(King_James)/Song_of_Solomon"),
	canonical(23, "Isaiah", "Bible_(King_James)/Isaiah"),
	canonical(24, "Jeremiah", "Bible_(King_James)/Jeremiah"),
	canonical(25, "Lamentations", "Bible_(King_James)/Lamentations"),
	canonical(26, "Ezekiel", "Bible_(King_James)/Ezekiel"),
	canonical(27, "Daniel", "Bible_(King_James)/Daniel"),
	canonical(28, "Hosea", "Bible_(King_James)/Hosea"),
	canonical(29, "Joel", "Bible_(King_James)/Joel"),
	canonical(30, "Amos", "Bible_(King_James)/Amos"),
	canonical(31, "Obadiah", "Bible_(King_James)/Obadiah"),
	canonical(32, "Jonah", "Bible_(King_James)/Jonah"),
	canonical(33, "Micah", "Bible_(King_James)/Micah"),
	canonical(34, "Nahum", "Bible_(King_James)/Nahum"),
	canonical(35, "Habakkuk", "Bible_(King_James)/Habakkuk"),
	canonical(36, "Zephaniah", "Bible_(King_James)/Zephaniah"),
	canonical(37, "Haggai", "Bible_(King_James)/Haggai"),
	canonical(38, "Zechariah", "Bible_(King_James)/Zechariah"),
	canonical(39, "Malachi", "Bible_(King_James)/Malachi"),
	// ---- New Testament ----
	canonical(40, "Matthew", "Bible_(King_James)/Matthew"),
	canonical(41, "Mark", "Bible_(King_James)/Mark"),
	canonical(42, "Luke", "Bible_(King_James)/Luke"),
	canonical(43, "John", "Bible_(King_James)/John"),
	canonical(44, "Acts", "Bible_(King_James)/Acts"),
	canonical(45, "Romans", "Bible_(King_James)/Romans"),
	canonical(46, "1 Corinthians", "Bible_(King_James)/I_Corinthians"),
	canonical(47, "2 Corinthians", "Bible_(King_James)/II_Corinthians"),
	canonical(48, "Galatians", "Bible_(King_James)/Galatians"),
	canonical(49, "Ephesians", "Bible_(King_James)/Ephesians"),
	canonical(50, "Philippians", "Bible_(King_James)/Philippians"),
	canonical(51, "Colossians", "Bible_(King_James)/Colossians"),
	canonical(52, "1 Thessalonians", "Bible_(King_James)/I_Thessalonians"),
	canonical(53, "2 Thessalonians", "Bible_(King_James)/II_Thessalonians"),
	canonical(54, "1 Timothy", "Bible_(King_James)/I_Timothy"),
	canonical(55, "2 Timothy", "Bible_(King_James)/II_Timothy"),
	canonical(56, "Titus", "Bible_(King_James)/Titus"),
	canonical(57, "Philemon", "Bible_(King_James)/Philemon"),
	canonical(58, "Hebrews", "Bible_(King_James)/Hebrews"),
	canonical(59, "James", "Bible_(King_James)/James"),
	canonical(60, "1 Peter", "Bible_(King_James)/I_Peter"),
	canonical(61, "2 Peter", "Bible_(King_James)/II_Peter"),
	canonical(62, "1 John", "Bible_(King_James)/I_John"),
	canonical(63, "2 John", "Bible_(King_James)/II_John"),
	canonical(64, "3 John", "Bible_(King_James)/III_John"),
	canonical(65, "Jude", "Bible_(King_James)/Jude"),
	canonical(66, "Revelation", "Bible_(King_James)/Revelation"),
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

// CanonicalBooks returns all 66 Protestant canonical books.
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
