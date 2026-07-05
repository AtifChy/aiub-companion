import { SearchInput } from "@/components/search-input";
import { Badge } from "@/components/ui/badge";
import { Button } from "@/components/ui/button";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { useDebounce } from "@/hooks/use-debounce";
import { logger } from "@/lib/logger";
import { cn } from "@/lib/utils";
import { Course, Service as RoutineService } from "@bindings/routine";
import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { Dialogs } from "@wailsio/runtime";
import {
  CalendarIcon,
  ClockIcon,
  FileSpreadsheetIcon,
  Loader2Icon,
  MapPinIcon,
  Trash2Icon,
  UserIcon,
} from "lucide-react";
import { useState } from "react";
import { toast } from "sonner";

type CourseStatus = "ongoing" | "upcoming" | "inactive";

const DAYS: readonly string[] = [
  "Sunday",
  "Monday",
  "Tuesday",
  "Wednesday",
  "Thursday",
  "Friday",
  "Saturday",
];

export default function RoutinePage() {
  const queryClient = useQueryClient();

  const { data: routine, isLoading: loading } = useQuery({
    queryKey: ["routine"],
    queryFn: () => RoutineService.GetUserRoutine(),
  });

  const [search, setSearch] = useState("");
  const searchDebounced = useDebounce(search, 300);

  const { data: searchResults, isLoading: isSearching } = useQuery({
    queryKey: ["courses", searchDebounced],
    queryFn: () => RoutineService.SearchOfferedCourses(searchDebounced),
    enabled: searchDebounced.trim().length > 0,
  });

  const addCourseMutation = useMutation({
    mutationFn: (classId: string) => RoutineService.AddToUserRoutine(classId),
    onSuccess: () => {
      setSearch("");
      void queryClient.invalidateQueries({ queryKey: ["routine"] });
      toast.success("Course added to routine");
    },
    onError: (err) => {
      logger.error("Failed to add course", err);
      toast.error("Failed to add course", {
        description: err.message,
      });
    },
  });

  const removeCourseMutation = useMutation({
    mutationFn: (classId: string) =>
      RoutineService.RemoveFromUserRoutine(classId),
    onSuccess: () => {
      toast.success("Course removed from routine");
      void queryClient.invalidateQueries({ queryKey: ["routine"] });
    },
    onError: (err) => {
      logger.error("Failed to remove course", err);
      toast.error("Failed to remove course", {
        description: err.message,
      });
    },
  });

  const importCoursesMutation = useMutation({
    mutationFn: (path: string) => RoutineService.ImportOfferedCourses(path),
    onSuccess: () => {
      void queryClient.invalidateQueries({ queryKey: ["courses"] });
      toast.success("Offered courses database updated");
    },
    onError: (err) => {
      logger.error("Failed to import courses", err);
      toast.error("Failed to import courses", {
        description: err.message,
      });
    },
  });

  const handleImportCourses = () => {
    Dialogs.OpenFile({
      Title: "Select Offered Courses Excel",
      Filters: [
        { DisplayName: "Excel files", Pattern: "*.xlsx;*.xls" },
        { DisplayName: "All files", Pattern: "*.*" },
      ],
    })
      .then((path) => {
        if (!path) return;
        importCoursesMutation.mutate(path);
      })
      .catch((err) => {
        logger.error("Failed to open file dialog", err);
        toast.error("Failed to open file dialog", {
          description: (err as Error).message,
        });
      });
  };

  const routineByDay = DAYS.reduce(
    (acc, day) => {
      acc[day] = routine?.filter((c) => c.day === day) ?? [];
      return acc;
    },
    {} as Record<string, Course[]>,
  );

  const stats: RoutineStatsProps["stats"] = {
    totalClasses: routine?.length ?? 0,
    studyHours:
      routine?.reduce((acc, c) => {
        const duration =
          parseTimeToMinutes(c.endTime) - parseTimeToMinutes(c.startTime);
        return acc + (duration > 0 ? duration / 60 : 0);
      }, 0) ?? 0,
    activeDays: Object.keys(routineByDay).filter((day) => {
      const dayCourses = routineByDay[day];
      return dayCourses !== undefined && dayCourses.length > 0;
    }).length,
  };

  return (
    <div className="animate-in fade-in-10 scrollbar-thumb-accent m-0.5 flex h-full scrollbar-thin scrollbar-gutter-both flex-col gap-6 overflow-auto p-6 duration-200 lg:p-10">
      <RoutineHeader onImport={handleImportCourses} />

      {routine && routine.length > 0 && !loading && (
        <RoutineStats stats={stats} />
      )}

      <CourseSearch
        search={search}
        setSearch={setSearch}
        searchResults={searchResults ?? []}
        isSearching={isSearching}
        onAddCourse={(classId) => addCourseMutation.mutate(classId)}
      />

      {loading ? (
        <div className="flex min-h-[300px] flex-1 items-center justify-center">
          <Loader2Icon className="text-muted-foreground size-8 animate-spin" />
        </div>
      ) : routine?.length === 0 ? (
        <EmptyState />
      ) : (
        <DayScheduleTimeline
          routineByDay={routineByDay}
          onRemoveCourse={(classId) => removeCourseMutation.mutate(classId)}
        />
      )}
    </div>
  );
}

interface RoutineHeaderProps {
  onImport: () => void;
}

function RoutineHeader({ onImport }: RoutineHeaderProps) {
  return (
    <div className="flex flex-col justify-between gap-4 sm:flex-row sm:items-center">
      <div>
        <h1 className="text-3xl font-bold tracking-tight">Class Routine</h1>
        <p className="text-muted-foreground mt-1 text-sm">
          Manage and visualize your weekly class schedule
        </p>
      </div>

      <Button variant="outline" className="w-fit" onClick={onImport}>
        <FileSpreadsheetIcon className="mr-2 size-4" />
        Import Offered Courses
      </Button>
    </div>
  );
}

interface RoutineStatsProps {
  stats: {
    totalClasses: number;
    studyHours: number;
    activeDays: number;
  };
}

function RoutineStats({ stats }: RoutineStatsProps) {
  return (
    <div className="animate-in fade-in mb-2 grid grid-cols-1 gap-4 duration-300 md:grid-cols-3">
      <Card className="bg-card/40 border-muted/30">
        <CardContent className="flex items-center gap-4">
          <div className="bg-primary/10 text-primary rounded-xl p-3">
            <ClockIcon className="size-5" />
          </div>
          <div>
            <p className="text-muted-foreground text-[10px] font-semibold tracking-wider uppercase">
              Weekly Study Load
            </p>
            <h3 className="text-lg font-bold">
              {stats.studyHours.toFixed(1)} Hours
            </h3>
          </div>
        </CardContent>
      </Card>

      <Card className="bg-card/40 border-muted/30">
        <CardContent className="flex items-center gap-4">
          <div className="rounded-xl bg-emerald-500/10 p-3 text-emerald-500">
            <CalendarIcon className="size-5" />
          </div>
          <div>
            <p className="text-muted-foreground text-[10px] font-semibold tracking-wider uppercase">
              Classes Scheduled
            </p>
            <h3 className="text-lg font-bold">
              {stats.totalClasses} Active Classes
            </h3>
          </div>
        </CardContent>
      </Card>

      <Card className="bg-card/40 border-muted/30">
        <CardContent className="flex items-center gap-4">
          <div className="rounded-xl bg-indigo-500/10 p-3 text-indigo-500">
            <MapPinIcon className="size-5" />
          </div>
          <div>
            <p className="text-muted-foreground text-[10px] font-semibold tracking-wider uppercase">
              Active Campus Days
            </p>
            <h3 className="text-lg font-bold">
              {stats.activeDays} Days / Week
            </h3>
          </div>
        </CardContent>
      </Card>
    </div>
  );
}

interface CourseSearchProps {
  search: string;
  setSearch: (s: string) => void;
  searchResults: Course[];
  isSearching: boolean;
  onAddCourse: (classId: string) => void;
}

function CourseSearch({
  search,
  setSearch,
  searchResults,
  isSearching,
  onAddCourse,
}: CourseSearchProps) {
  return (
    <div className="relative z-50">
      <SearchInput
        value={search}
        onValueChange={setSearch}
        placeholder="Search and add offered courses (e.g., Computer Networks, CSE 3101)..."
        searching={isSearching}
        showClearButton={false}
        className="[&>div:last-child]:right-3 [&>input]:h-10 [&>input]:pl-9 [&>input]:text-sm [&>input]:shadow-sm [&>svg]:size-4 [&>svg:first-child]:left-3 [&>svg:last-child]:right-3 [&>svg:last-child]:size-5"
      />

      {search !== "" && searchResults.length > 0 && (
        <Card className="scrollbar-thumb-accent bg-popover/95 border-muted/40 fade-in-10 animate-in absolute mt-2 max-h-80 w-full scrollbar-thin overflow-y-auto rounded-lg border pt-0 pb-0 shadow-xl backdrop-blur-md duration-200">
          <div className="flex flex-col">
            {searchResults.map((course) => (
              <button
                type="button"
                key={course.classID}
                tabIndex={0}
                onClick={() => onAddCourse(course.classID)}
                onKeyDown={(e) => {
                  if (e.key === "Enter") {
                    e.preventDefault();
                    onAddCourse(course.classID);
                  }
                }}
                className={cn(
                  "hover:bg-accent/60 flex cursor-pointer items-center justify-between border-b px-4 py-3 transition-colors last:border-0",
                  "focus-visible:ring-ring/50 outline-none focus-visible:rounded focus-visible:ring-3 focus-visible:ring-inset",
                )}
              >
                <div className="flex flex-col gap-1">
                  <div className="flex items-center gap-2 text-sm font-semibold">
                    {course.courseCode && (
                      <Badge
                        variant="secondary"
                        className="px-1.5 py-0.5 font-mono text-[9px]"
                      >
                        {course.courseCode}
                      </Badge>
                    )}
                    <span>{course.courseTitle}</span>
                    <span className="text-muted-foreground text-xs font-normal">
                      [{course.section}]
                    </span>
                  </div>
                  <div className="text-muted-foreground flex items-center gap-4 text-[11px]">
                    <span className="flex items-center gap-1">
                      <UserIcon className="size-3" /> {course.faculty}
                    </span>
                    <span className="flex items-center gap-1">
                      <CalendarIcon className="size-3" /> {course.day}
                    </span>
                    <span className="flex items-center gap-1">
                      <ClockIcon className="size-3" /> {course.startTime} -{" "}
                      {course.endTime}
                    </span>
                  </div>
                </div>
              </button>
            ))}
          </div>
        </Card>
      )}

      {search !== "" && !isSearching && searchResults.length === 0 && (
        <Card className="text-muted-foreground animate-in fade-in-10 border-muted/40 absolute mt-2 w-full rounded-lg border p-4 text-center text-xs shadow-lg duration-200">
          No matching offered courses found in the database.
        </Card>
      )}
    </div>
  );
}

function EmptyState() {
  return (
    <div className="flex min-h-[300px] flex-1 flex-col items-center justify-center py-20 text-center">
      <div className="bg-muted mb-4 rounded-full p-4">
        <CalendarIcon className="text-muted-foreground/60 size-10" />
      </div>
      <h3 className="text-lg font-bold">Your routine is empty</h3>
      <p className="text-muted-foreground mt-1 max-w-sm text-xs leading-relaxed">
        Search for your offered courses in the bar above or import an offered
        Excel sheet to construct your weekly timeline.
      </p>
    </div>
  );
}

interface DayScheduleTimelineProps {
  routineByDay: Record<string, Course[]>;
  onRemoveCourse: (classId: string) => void;
}

function DayScheduleTimeline({
  routineByDay,
  onRemoveCourse,
}: DayScheduleTimelineProps) {
  return (
    <div className="flex flex-col gap-8 pb-12">
      {DAYS.map((day) => {
        const dayCourses = routineByDay[day];
        if (!dayCourses || dayCourses.length === 0) return null;

        return (
          <div
            key={day}
            className="animate-in fade-in flex flex-col gap-6 border-b pb-6 duration-300 last:border-0 lg:flex-row"
          >
            {/* Sticky Side Day Column */}
            <div className="flex shrink-0 items-baseline justify-between gap-1 lg:w-36 lg:flex-col lg:items-start">
              <h2 className="text-foreground text-lg font-black tracking-tight">
                {day}
              </h2>
              <Badge
                variant="secondary"
                className="px-2 py-0.5 text-[10px] font-bold"
              >
                {dayCourses.length}{" "}
                {dayCourses.length === 1 ? "Class" : "Classes"}
              </Badge>
            </div>

            {/* Day Grid Courses */}
            <div className="grid flex-1 grid-cols-1 gap-4 md:grid-cols-2 xl:grid-cols-3">
              {dayCourses.map((course) => (
                <CourseCard
                  key={course.classID}
                  course={course}
                  onRemoveCourse={onRemoveCourse}
                />
              ))}
            </div>
          </div>
        );
      })}
    </div>
  );
}

interface CourseCardProps {
  course: Course;
  onRemoveCourse: (classId: string) => void;
}

function CourseCard({ course, onRemoveCourse }: CourseCardProps) {
  const status = getCourseStatus(course);
  const isLab = course.type?.toLowerCase().includes("lab");

  return (
    <Card
      className={cn(
        "group relative overflow-hidden border-l-4 transition-all duration-300 hover:-translate-y-1 hover:shadow-md",
        isLab ? "border-l-emerald-500" : "border-l-indigo-500",
        status === "ongoing" &&
          "border-l-emerald-500 bg-emerald-500/5 shadow-md ring-2 shadow-emerald-500/5 ring-emerald-500",
      )}
    >
      {/* Live Status Tracker Pulses */}
      {status === "ongoing" && (
        <span className="absolute top-4 right-4 flex h-2 w-2">
          <span className="absolute inline-flex h-full w-full animate-ping rounded-full bg-emerald-400 opacity-75"></span>
          <span className="relative inline-flex h-2 w-2 rounded-full bg-emerald-500"></span>
        </span>
      )}

      <CardHeader>
        <div className="mb-1 flex flex-wrap items-center gap-1.5">
          {course.courseCode && (
            <Badge
              variant="secondary"
              className="px-1 py-0 font-mono text-[9px] opacity-80"
            >
              {course.courseCode}
            </Badge>
          )}
          <Badge
            variant="outline"
            className="-ml-1 px-2 py-0 text-[0.65rem] font-medium"
          >
            {course.type}
          </Badge>
          {status === "ongoing" && (
            <Badge className="border-0 bg-emerald-500 px-1.5 py-0 text-[8px] font-bold tracking-wider text-white uppercase">
              Ongoing
            </Badge>
          )}
          {status === "upcoming" && (
            <Badge className="border-0 bg-amber-500 px-1.5 py-0 text-[8px] font-bold tracking-wider text-white uppercase">
              Up Next
            </Badge>
          )}
        </div>
        <CardTitle className="pr-6 text-sm leading-snug font-bold">
          {course.courseTitle}{" "}
          <span className="text-muted-foreground mt-0.5 block text-xs font-normal">
            Section {course.section}
          </span>
        </CardTitle>
      </CardHeader>

      <CardContent className="flex flex-col gap-2 pt-0">
        <div className="text-muted-foreground flex items-center gap-2 text-xs">
          <ClockIcon className="text-foreground/60 size-3.5" />
          <span className="text-foreground/90 font-medium">
            {course.startTime} - {course.endTime}
          </span>
        </div>
        <div className="text-muted-foreground flex items-center gap-2 text-xs">
          <MapPinIcon className="text-foreground/60 size-3.5" />
          <span className="truncate">
            Room {course.room} ({course.department})
          </span>
        </div>
        <div className="text-muted-foreground flex items-center gap-2 text-xs">
          <UserIcon className="text-foreground/60 size-3.5" />
          <span className="truncate">{course.faculty}</span>
        </div>
      </CardContent>

      {/* Interactive Trash Trigger */}
      <Button
        variant="destructive"
        size="icon"
        className="absolute right-3 bottom-3 h-8 w-8 rounded-md opacity-0 transition-opacity duration-200 group-hover:opacity-100"
        onClick={() => void onRemoveCourse(course.classID)}
      >
        <Trash2Icon className="size-4" />
      </Button>
    </Card>
  );
}

// Helper to determine if a course is ongoing, upcoming (starts within 60 minutes), or inactive
const getCourseStatus = (course: Course): CourseStatus => {
  const now = new Date();
  const daysMap = [
    "Sunday",
    "Monday",
    "Tuesday",
    "Wednesday",
    "Thursday",
    "Friday",
    "Saturday",
  ];
  const currentDay = daysMap[now.getDay()];

  if (course.day !== currentDay) return "inactive";

  const currentMinutes = now.getHours() * 60 + now.getMinutes();
  const startMinutes = parseTimeToMinutes(course.startTime);
  const endMinutes = parseTimeToMinutes(course.endTime);

  if (currentMinutes >= startMinutes && currentMinutes <= endMinutes) {
    return "ongoing";
  } else if (
    currentMinutes < startMinutes &&
    startMinutes - currentMinutes <= 60
  ) {
    return "upcoming";
  }
  return "inactive";
};

// Helper to parse time string ("08:00 AM") to minutes since midnight
const parseTimeToMinutes = (timeStr: string | undefined): number => {
  if (!timeStr) return 0;
  const parts = timeStr.trim().split(" ");
  if (parts.length < 2) return 0;
  const timeVal = parts[0];
  const modifier = parts[1];
  if (!timeVal || !modifier) return 0;

  const [hoursStr, minutesStr] = timeVal.split(":");
  let hours = parseInt(hoursStr || "0", 10);
  const minutes = parseInt(minutesStr || "0", 10);

  if (modifier === "PM" && hours < 12) hours += 12;
  if (modifier === "AM" && hours === 12) hours = 0;

  return hours * 60 + minutes;
};
