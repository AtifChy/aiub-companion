import { Badge } from "@/components/ui/badge";
import { Button } from "@/components/ui/button";
import { Card, CardContent } from "@/components/ui/card";
import { Input } from "@/components/ui/input";
import { Separator } from "@/components/ui/separator";
import { Browser } from "@wailsio/runtime";
import {
  ChevronDownIcon,
  ChevronUpIcon,
  ExternalLinkIcon,
  GitBranchIcon,
  HelpCircleIcon,
  InfoIcon,
  SearchIcon,
} from "lucide-react";
import { useState } from "react";
import { toast } from "sonner";

interface FAQItem {
  question: string;
  answer: string;
  category: "general" | "routine" | "notices" | "settings";
}

const faqs: FAQItem[] = [
  {
    category: "general",
    question: "What is AIUB Companion?",
    answer:
      "It is an open-source productivity app designed for American International University-Bangladesh (AIUB) students. It aggregates portal notices, displays real-time countdowns, calculates CGPA, and manages class schedules.",
  },
  {
    category: "routine",
    question: "How do I import offered courses?",
    answer:
      "Go to the Routine page and click 'Import Offered Courses'. Select the Excel spreadsheet file (.xlsx) downloaded from the official AIUB portal. The application will scan and map all class sessions automatically.",
  },
  {
    category: "routine",
    question: "Why is a class missing from my search?",
    answer:
      "Ensure that you have imported the latest offered course Excel spreadsheet. If the file is imported successfully and the course is still missing, check if the sheet columns match the standard portal template.",
  },
  {
    category: "notices",
    question: "How frequently does the notice scraper update?",
    answer:
      "By default, notices are scraped and synchronized every 30 minutes. You can adjust this duration or trigger manual synchronization in the Settings panel.",
  },
  {
    category: "notices",
    question: "Can I view downloaded notice attachments offline?",
    answer:
      "Yes. Once notice details and PDF attachments are fetched/opened, they are cached locally on your device for offline reading.",
  },
  {
    category: "settings",
    question: "Where are my personal settings and data stored?",
    answer:
      "All configuration settings and cached notices are stored locally inside your user profile configuration directory (AppData/Roaming on Windows) under 'aiub-companion'. No personal data is sent to external servers.",
  },
] as const;

export default function HelpPage() {
  const [searchQuery, setSearchQuery] = useState("");

  return (
    <div className="bg-background animate-in fade-in-10 scrollbar-thumb-accent h-full scrollbar-thin scrollbar-gutter-both space-y-8 overflow-y-auto p-6 duration-200 lg:p-10">
      <Header />

      {/* FAQ Section */}
      <div className="space-y-4">
        <div className="flex flex-col justify-between gap-4 sm:flex-row sm:items-center">
          <h2 className="flex items-center gap-2 text-xl font-bold">
            <HelpCircleIcon className="text-muted-foreground h-5 w-5" />
            Frequently Asked Questions
          </h2>

          {/* Search Box */}
          <div className="relative w-full sm:w-80">
            <SearchIcon className="text-muted-foreground absolute top-1/2 left-3 h-4 w-4 -translate-y-1/2" />
            <Input
              type="text"
              placeholder="Search help topics..."
              value={searchQuery}
              onChange={(e) => setSearchQuery(e.target.value)}
              className="pl-9"
            />
          </div>
        </div>

        <FAQList searchQuery={searchQuery} />
      </div>

      <Separator />
      <Footer />
    </div>
  );
}

function Header() {
  return (
    <div className="space-y-2">
      <h1 className="text-3xl font-bold tracking-tight">Help</h1>
      <p className="text-muted-foreground text-sm">
        Find answers, view quick guides, and get the most out of your AIUB
        Companion.
      </p>
    </div>
  );
}

function FAQList({ searchQuery }: { searchQuery: string }) {
  const [expandedIndex, setExpandedIndex] = useState<number | null>(null);

  const filteredFaqs = faqs.filter(
    (faq) =>
      faq.question.toLowerCase().includes(searchQuery.toLowerCase()) ||
      faq.answer.toLowerCase().includes(searchQuery.toLowerCase()) ||
      faq.category.toLowerCase().includes(searchQuery.toLowerCase()),
  );

  const toggleExpand = (index: number) => {
    setExpandedIndex(expandedIndex === index ? null : index);
  };

  return (
    <div className="bg-card divide-y overflow-hidden rounded-lg border">
      {filteredFaqs.length > 0 ? (
        filteredFaqs.map((faq, idx) => (
          <FAQItem
            key={idx}
            faq={faq}
            expanded={expandedIndex === idx}
            onExpand={() => toggleExpand(idx)}
          />
        ))
      ) : (
        <div className="text-muted-foreground p-8 text-center">
          No results found matching "{searchQuery}"
        </div>
      )}
    </div>
  );
}

function FAQItem({
  faq,
  expanded,
  onExpand,
}: {
  faq: FAQItem;
  expanded: boolean;
  onExpand: () => void;
}) {
  return (
    <div className="hover:bg-muted/30 transition-colors">
      <button
        onClick={onExpand}
        className="flex w-full items-center justify-between p-4 text-left font-medium outline-none"
      >
        <div className="flex items-center gap-3">
          {getCategoryBadge(faq.category)}
          <span className="text-sm sm:text-base">{faq.question}</span>
        </div>
        {expanded ? (
          <ChevronUpIcon className="text-muted-foreground size-4 shrink-0" />
        ) : (
          <ChevronDownIcon className="text-muted-foreground size-4 shrink-0" />
        )}
      </button>
      {expanded && (
        <div className="text-muted-foreground animate-in fade-in slide-in-from-top-1 px-4 pb-4 text-sm leading-relaxed duration-200">
          {faq.answer}
        </div>
      )}
    </div>
  );
}

function Footer() {
  return (
    <div className="flex flex-col gap-6 md:flex-row md:justify-between">
      {/* Support Links */}
      <div className="space-y-3">
        <h3 className="text-lg font-semibold">Useful Resources</h3>
        <div className="flex flex-wrap gap-3">
          <Button
            variant="outline"
            className="gap-2"
            onClick={() => {
              Browser.OpenURL(
                "https://github.com/atifchy/aiub-companion",
              ).catch((err) =>
                toast.error(`Failed to open URL`, {
                  description: err instanceof Error ? err.message : String(err),
                }),
              );
            }}
          >
            <GitBranchIcon className="size-4" />
            GitHub Repository
            <ExternalLinkIcon className="size-3" />
          </Button>
          <Button
            variant="outline"
            className="gap-2"
            onClick={() => {
              Browser.OpenURL("https://www.aiub.edu").catch((err) =>
                toast.error(`Failed to open URL`, {
                  description: err instanceof Error ? err.message : String(err),
                }),
              );
            }}
          >
            AIUB Website
            <ExternalLinkIcon className="h-3 w-3" />
          </Button>
        </div>
      </div>

      {/* Diagnostics Card */}
      <Card className="bg-muted/40">
        <CardContent className="flex gap-3">
          <InfoIcon className="text-muted-foreground mt-0.5 h-5 w-5 shrink-0" />
          <div className="text-muted-foreground space-y-1 text-xs">
            <p className="text-foreground font-semibold">
              App Diagnostics Information
            </p>
            <p>Storage: Local SQLite / JSON files</p>
            <p>Configuration Path: AppData/Roaming/aiub-companion</p>
            <p>Logs level: Configurable via Settings page</p>
          </div>
        </CardContent>
      </Card>
    </div>
  );
}

const getCategoryBadge = (category: string) => {
  switch (category) {
    case "routine":
      return (
        <Badge
          variant="secondary"
          className="border-none bg-blue-500/10 text-blue-500"
        >
          Routine
        </Badge>
      );
    case "notices":
      return (
        <Badge
          variant="secondary"
          className="border-none bg-emerald-500/10 text-emerald-500"
        >
          Notices
        </Badge>
      );
    case "settings":
      return (
        <Badge
          variant="secondary"
          className="border-none bg-amber-500/10 text-amber-500"
        >
          Settings
        </Badge>
      );
    default:
      return (
        <Badge
          variant="secondary"
          className="border-none bg-purple-500/10 text-purple-500"
        >
          General
        </Badge>
      );
  }
};
