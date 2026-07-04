import { Badge } from "@/components/ui/badge";
import { Button } from "@/components/ui/button";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { Progress } from "@/components/ui/progress";
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";
import { useDelayedLoading } from "@/hooks/use-delayed-loading";
import { logger } from "@/lib/logger";
import { cn } from "@/lib/utils";
import {
  Service as CalendarService,
  CalendarType,
  type AcademicEvent,
} from "@bindings/calendar";
import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import {
  BookOpenIcon,
  CalendarDaysIcon,
  ClockIcon,
  GraduationCapIcon,
  Loader2Icon,
  RefreshCwIcon,
  TargetIcon,
  TrophyIcon,
} from "lucide-react";
import { useState } from "react";
import { toast } from "sonner";

const CALENDAR_TYPES = [
  { value: CalendarType.CalendarStandard, label: "Standard" },
  { value: CalendarType.CalendarLLBBPharm, label: "LLB & BPharm" },
];

const CATEGORY_STYLES: Record<string, string> = {
  exam: "border-red-500/20 bg-red-500/10 text-red-500",
  payment: "border-green-500/20 bg-green-500/10 text-green-500",
  registration: "border-violet-500/20 bg-violet-500/10 text-violet-500",
  academic: "border-blue-500/20 bg-blue-500/10 text-blue-500",
  deadline: "border-orange-500/20 bg-orange-500/10 text-orange-500",
  lab: "border-cyan-500/20 bg-cyan-500/10 text-cyan-500",
  break: "border-yellow-500/20 bg-yellow-500/10 text-yellow-500",
};

const CATEGORY_ICONS: Record<string, React.ElementType> = {
  exam: TargetIcon,
  payment: BookOpenIcon,
  registration: GraduationCapIcon,
  academic: BookOpenIcon,
  deadline: ClockIcon,
  lab: BookOpenIcon,
  break: CalendarDaysIcon,
};

export default function SemesterPage() {
  const queryClient = useQueryClient();

  const [calendarType, setCalendarType] = useState<CalendarType>(
    CalendarType.CalendarStandard,
  );
  const [view, setView] = useState<"timeline" | "all">("timeline");

  const calendarQuery = useQuery({
    queryKey: ["calendar", calendarType],
    queryFn: async () => {
      const [calendar, currentWeek, nextExam, upcomingEvents] =
        await Promise.all([
          CalendarService.GetAcademicCalendar(calendarType),
          CalendarService.GetCurrentWeek(calendarType),
          CalendarService.GetNextExam(calendarType),
          CalendarService.GetUpcomingEvents(calendarType, 10),
        ]);
      return { calendar, currentWeek, nextExam, upcomingEvents };
    },
  });

  const { data, isLoading } = calendarQuery;
  const calendar = data?.calendar ?? null;
  const currentWeek = data?.currentWeek ?? 0;
  const nextExam = data?.nextExam ?? null;
  const upcomingEvents = data?.upcomingEvents ?? [];

  const { mutate, isPending } = useMutation({
    mutationFn: (type: CalendarType) => CalendarService.RefreshCalendar(type),
    onSuccess: (_, type) => {
      void queryClient.invalidateQueries({
        queryKey: ["calendar", type],
      });
      toast.success("Calendar refreshed");
    },
    onError: (err) => {
      logger.error("Failed to refresh calendar", err);
      toast.error("Failed to refresh calendar", {
        description: err.message,
      });
    },
  });

  const showLoading = useDelayedLoading(isLoading);

  if (isLoading) {
    if (!showLoading) return null;
    return (
      <div className="flex h-full items-center justify-center">
        <Loader2Icon className="text-muted-foreground h-8 w-8 animate-spin" />
      </div>
    );
  }

  const progress = calendar?.totalWeeks
    ? (currentWeek / calendar.totalWeeks) * 100
    : 0;

  return (
    <div className="scrollbar-thumb-accent animate-in fade-in-10 m-0.5 h-full scrollbar-thin scrollbar-gutter-both overflow-auto duration-200">
      <div className="space-y-6 p-6 lg:p-10">
        {/* Header */}
        <div className="flex flex-col gap-4 sm:flex-row sm:items-center sm:justify-between">
          <div>
            <h1 className="text-3xl font-bold">Semester Dashboard</h1>
            {calendar?.semester && (
              <p className="text-muted-foreground text-sm capitalize">
                {calendar.semester}
              </p>
            )}
          </div>
          <div className="flex items-center gap-2">
            <Select
              items={CALENDAR_TYPES}
              value={calendarType}
              onValueChange={(v) => v && setCalendarType(v)}
            >
              <SelectTrigger className="w-[180px]">
                <SelectValue placeholder="Calendar type" />
              </SelectTrigger>
              <SelectContent className="p-1">
                {CALENDAR_TYPES.map((ct) => (
                  <SelectItem key={ct.value} value={ct.value}>
                    {ct.label}
                  </SelectItem>
                ))}
              </SelectContent>
            </Select>
            <Button
              variant="outline"
              size="icon"
              onClick={() => mutate(calendarType)}
              disabled={isPending}
            >
              <RefreshCwIcon
                className={cn("h-4 w-4", isPending && "animate-spin")}
              />
            </Button>
          </div>
        </div>

        {/* Stats Grid */}
        <div className="grid gap-4 sm:grid-cols-2 lg:grid-cols-4">
          <Card>
            <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
              <CardTitle className="text-sm font-medium">
                Current Week
              </CardTitle>
              <CalendarDaysIcon className="text-muted-foreground h-4 w-4" />
            </CardHeader>
            <CardContent>
              <div className="text-2xl font-bold">
                Week {currentWeek || "—"}
              </div>
              {calendar?.totalWeeks && (
                <p className="text-muted-foreground text-xs">
                  of {calendar.totalWeeks} total weeks
                </p>
              )}
            </CardContent>
          </Card>

          <Card>
            <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
              <CardTitle className="text-sm font-medium">
                Semester Progress
              </CardTitle>
              <TrophyIcon className="text-muted-foreground h-4 w-4" />
            </CardHeader>
            <CardContent>
              <div className="text-2xl font-bold">{progress.toFixed(0)}%</div>
              <Progress value={progress} className="*:bg-accent mt-2 h-2" />
            </CardContent>
          </Card>

          <Card>
            <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
              <CardTitle className="text-sm font-medium">Next Exam</CardTitle>
              <TargetIcon className="text-muted-foreground h-4 w-4" />
            </CardHeader>
            <CardContent>
              <div className="text-2xl font-bold">
                {nextExam ? nextExam.title : "—"}
              </div>
              {nextExam?.date && (
                <p className="text-muted-foreground text-xs">
                  {new Date(nextExam.date).toLocaleDateString()}
                </p>
              )}
            </CardContent>
          </Card>

          <Card>
            <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
              <CardTitle className="text-sm font-medium">
                Total Events
              </CardTitle>
              <BookOpenIcon className="text-muted-foreground h-4 w-4" />
            </CardHeader>
            <CardContent>
              <div className="text-2xl font-bold">
                {calendar?.events?.length ?? 0}
              </div>
              <p className="text-muted-foreground text-xs">
                events this semester
              </p>
            </CardContent>
          </Card>
        </div>

        {/* View Toggle */}
        <div className="flex gap-2">
          <Button
            variant={view === "timeline" ? "default" : "outline"}
            size="sm"
            onClick={() => setView("timeline")}
          >
            Timeline
          </Button>
          <Button
            variant={view === "all" ? "default" : "outline"}
            size="sm"
            onClick={() => setView("all")}
          >
            All Events
          </Button>
        </div>

        {/* Events Content */}
        <Card>
          <CardHeader>
            <CardTitle>
              {view === "timeline"
                ? "Upcoming Events"
                : `All Events (${calendar?.events?.length ?? 0})`}
            </CardTitle>
          </CardHeader>
          <CardContent>
            {view === "timeline" ? (
              <div className="relative space-y-4">
                {/* Timeline line */}
                <div className="bg-border absolute top-2 bottom-2 left-[11px] w-px" />

                {upcomingEvents.length === 0 ? (
                  <p className="text-muted-foreground py-8 text-center">
                    No upcoming events
                  </p>
                ) : (
                  upcomingEvents.map((event, idx) => (
                    <EventItem key={idx} event={event} />
                  ))
                )}
              </div>
            ) : (
              <div className="space-y-2">
                {calendar?.events?.map((event, idx) => (
                  <div
                    key={idx}
                    className="flex items-center justify-between rounded-lg border p-3"
                  >
                    <div className="space-y-1">
                      <p className="font-medium">{event.title}</p>
                      <p className="text-muted-foreground text-sm">
                        {formatEventDate(event)}
                      </p>
                    </div>
                    <Badge
                      variant="outline"
                      className={cn(
                        "capitalize",
                        CATEGORY_STYLES[event.category] ??
                          "border-gray-500/20 bg-gray-500/10 text-gray-500",
                      )}
                    >
                      {event.category}
                    </Badge>
                  </div>
                ))}
              </div>
            )}
          </CardContent>
        </Card>
      </div>
    </div>
  );
}

function EventItem({ event }: { event: AcademicEvent }) {
  const Icon = CATEGORY_ICONS[event.category] ?? BookOpenIcon;
  const style =
    CATEGORY_STYLES[event.category] ??
    "border-gray-500/20 bg-gray-500/10 text-gray-500";

  return (
    <div className="relative flex gap-4 pl-8">
      {/* Timeline dot */}
      <div
        className={cn(
          "absolute top-1 left-0 flex h-6 w-6 items-center justify-center rounded-full border",
          style,
        )}
      >
        <Icon className="h-3 w-3" />
      </div>

      <div className="flex-1 space-y-1 rounded-lg border p-3">
        <div className="flex items-start justify-between gap-2">
          <div className="space-y-1">
            <p className="leading-none font-medium">{event.title}</p>
            <p className="text-muted-foreground text-sm">
              {formatEventDate(event)}
              {event.week && ` • Week ${event.week}`}
            </p>
          </div>
          <Badge variant="outline" className={cn("capitalize", style)}>
            {event.category}
          </Badge>
        </div>
      </div>
    </div>
  );
}

function formatEventDate(event: AcademicEvent): string {
  const date = event.date ? new Date(event.date) : null;
  const endDate = event.endDate ? new Date(event.endDate) : null;

  if (!date) return "TBD";

  const opts: Intl.DateTimeFormatOptions = {
    month: "short",
    day: "numeric",
  };

  if (date.getFullYear() !== new Date().getFullYear()) {
    opts.year = "numeric";
  }

  if (endDate) {
    return `${date.toLocaleDateString(undefined, opts)} – ${endDate.toLocaleDateString(undefined, opts)}`;
  }

  return date.toLocaleDateString(undefined, opts);
}
