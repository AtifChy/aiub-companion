import { Browser } from "@wailsio/runtime";
import {
  ChevronDownIcon,
  ChevronUpIcon,
  ExternalLinkIcon,
  GitBranchIcon,
  HelpCircleIcon,
} from "lucide-react";
import { useState } from "react";
import { toast } from "sonner";

import { SearchInput } from "@/components/search-input";
import { Badge } from "@/components/ui/badge";
import { Button } from "@/components/ui/button";
import { Separator } from "@/components/ui/separator";

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
  return (
    <div className="mr-0.5 h-full animate-in scrollbar-thin scrollbar-thumb-accent scrollbar-gutter-both space-y-8 overflow-y-auto bg-background p-6 duration-200 fade-in-10 lg:p-10">
      <Header />
      <FAQ />
      <Separator />
      <Footer />
    </div>
  );
}

function Header() {
  return (
    <div className="space-y-2">
      <h1 className="text-3xl font-bold tracking-tight">Help</h1>
      <p className="text-sm text-muted-foreground">
        Find answers, view quick guides, and get the most out of your AIUB Companion.
      </p>
    </div>
  );
}

function FAQ() {
  const [search, setSearch] = useState("");
  return (
    <div className="space-y-4">
      <div className="flex flex-col justify-between gap-4 sm:flex-row sm:items-center">
        <h2 className="flex items-center gap-2 text-xl font-bold">
          <HelpCircleIcon className="h-5 w-5 text-muted-foreground" />
          Frequently Asked Questions
        </h2>

        <SearchInput
          value={search}
          onValueChange={setSearch}
          placeholder="Search help topics..."
          className="w-full sm:w-80"
        />
      </div>

      <FAQList searchQuery={search} />
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
    <div className="divide-y overflow-hidden rounded-lg border bg-card">
      {filteredFaqs.length > 0 ? (
        filteredFaqs.map((faq, idx) => (
          <FAQItem
            key={faq.category + String(idx)}
            faq={faq}
            expanded={expandedIndex === idx}
            onExpand={() => {
              toggleExpand(idx);
            }}
          />
        ))
      ) : (
        <div className="p-8 text-center text-muted-foreground">
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
    <div className="transition-colors hover:bg-muted/30">
      <button
        type="button"
        onClick={onExpand}
        className="flex w-full items-center justify-between p-4 text-left font-medium outline-none"
      >
        <div className="flex items-center gap-3">
          {getCategoryBadge(faq.category)}
          <span className="text-sm sm:text-base">{faq.question}</span>
        </div>
        {expanded ? (
          <ChevronUpIcon className="size-4 shrink-0 text-muted-foreground" />
        ) : (
          <ChevronDownIcon className="size-4 shrink-0 text-muted-foreground" />
        )}
      </button>
      {expanded && (
        <div className="animate-in px-4 pb-4 text-sm leading-relaxed text-muted-foreground duration-200 fade-in slide-in-from-top-1">
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
              Browser.OpenURL("https://github.com/atifchy/aiub-companion").catch((err: unknown) =>
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
              Browser.OpenURL("https://www.aiub.edu").catch((err: unknown) =>
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
    </div>
  );
}

const getCategoryBadge = (category: string) => {
  switch (category) {
    case "routine":
      return (
        <Badge variant="secondary" className="border-none bg-blue-500/10 text-blue-500">
          Routine
        </Badge>
      );
    case "notices":
      return (
        <Badge variant="secondary" className="border-none bg-emerald-500/10 text-emerald-500">
          Notices
        </Badge>
      );
    case "settings":
      return (
        <Badge variant="secondary" className="border-none bg-amber-500/10 text-amber-500">
          Settings
        </Badge>
      );
    default:
      return (
        <Badge variant="secondary" className="border-none bg-purple-500/10 text-purple-500">
          General
        </Badge>
      );
  }
};
