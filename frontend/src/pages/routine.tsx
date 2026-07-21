import { Course, Service as RoutineService, Schedule } from "@bindings/routine";
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

import { SearchInput } from "@/components/search-input";
import { Badge } from "@/components/ui/badge";
import { Button } from "@/components/ui/button";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { useDebounce } from "@/hooks/use-debounce";
import { logger } from "@/lib/logger";
import { DAYS, getCourseStatus, parseTimeToMinutes } from "@/lib/routine";
import { cn } from "@/lib/utils";

export default function RoutinePage() {
  const queryClient = useQueryClient();

  const { data: routine, isLoading: loading } = useQuery({
    queryKey: ["routine"],
    queryFn: () => RoutineService.ListUserRoutine(),
  });

  const { mutate: removeCourse } = useMutation({
    mutationFn: (classId: string) => RoutineService.RemoveFromUserRoutine(classId),
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

  const { mutate: importCourses } = useMutation({
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
        importCourses(path);
      })
      .catch((err: unknown) => {
        toast.error("Failed to open file dialog", {
          description: (err as Error).message,
        });
      });
  };

  const routineByDay = DAYS.reduce<Record<string, Course[]>>((acc, day) => {
    acc[day] = routine?.filter((course) => course.schedules.some((s) => s.day === day)) ?? [];
    return acc;
  }, {});

  const stats: RoutineStatsProps["stats"] = {
    totalClasses: routine?.length ?? 0,
    studyHours:
      routine?.reduce((total, course) => {
        return (
          total +
          course.schedules.reduce((sum, schedule) => {
            const duration =
              parseTimeToMinutes(schedule.endTime) - parseTimeToMinutes(schedule.startTime);
            return sum + Math.max(duration, 0) / 60;
          }, 0)
        );
      }, 0) ?? 0,
    activeDays: Object.keys(routineByDay).filter((day) => {
      const dayCourses = routineByDay[day];
      return dayCourses !== undefined && dayCourses.length > 0;
    }).length,
  };

  return (
    <div className="mr-0.5 flex h-full animate-in scrollbar-thin scrollbar-thumb-accent scrollbar-gutter-both flex-col gap-6 overflow-auto p-6 duration-200 fade-in-10 lg:p-10">
      <RoutineHeader onImport={handleImportCourses} />

      {routine && routine.length > 0 && <RoutineStats stats={stats} />}

      <CourseSearch />

      {loading ? (
        <div className="flex min-h-[300px] flex-1 items-center justify-center">
          <Loader2Icon className="size-8 animate-spin text-muted-foreground" />
        </div>
      ) : routine?.length === 0 ? (
        <EmptyState />
      ) : (
        <DayScheduleTimeline
          routineByDay={routineByDay}
          onRemoveCourse={(classId) => removeCourse(classId)}
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
        <p className="mt-1 text-sm text-muted-foreground">
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
    <div className="mb-2 grid animate-in grid-cols-1 gap-4 duration-300 fade-in md:grid-cols-3">
      <Card className="border-muted/30 bg-card/40">
        <CardContent className="flex items-center gap-4">
          <div className="rounded-xl bg-primary/10 p-3 text-primary">
            <ClockIcon className="size-5" />
          </div>
          <div>
            <p className="text-[10px] font-semibold tracking-wider text-muted-foreground uppercase">
              Weekly Study Load
            </p>
            <h3 className="text-lg font-bold">{stats.studyHours.toFixed(1)} Hours</h3>
          </div>
        </CardContent>
      </Card>

      <Card className="border-muted/30 bg-card/40">
        <CardContent className="flex items-center gap-4">
          <div className="rounded-xl bg-emerald-500/10 p-3 text-emerald-500">
            <CalendarIcon className="size-5" />
          </div>
          <div>
            <p className="text-[10px] font-semibold tracking-wider text-muted-foreground uppercase">
              Classes Scheduled
            </p>
            <h3 className="text-lg font-bold">{stats.totalClasses} Active Classes</h3>
          </div>
        </CardContent>
      </Card>

      <Card className="border-muted/30 bg-card/40">
        <CardContent className="flex items-center gap-4">
          <div className="rounded-xl bg-indigo-500/10 p-3 text-indigo-500">
            <MapPinIcon className="size-5" />
          </div>
          <div>
            <p className="text-[10px] font-semibold tracking-wider text-muted-foreground uppercase">
              Active Campus Days
            </p>
            <h3 className="text-lg font-bold">{stats.activeDays} Days / Week</h3>
          </div>
        </CardContent>
      </Card>
    </div>
  );
}

function CourseSearch() {
  const queryClient = useQueryClient();

  const [search, setSearch] = useState("");
  const searchDebounced = useDebounce(search, 300);

  const [isFocused, setIsFocused] = useState(false);

  const { data: searchResults, isLoading: isSearching } = useQuery({
    queryKey: ["courses", searchDebounced],
    queryFn: () => RoutineService.SearchOfferedCourses(searchDebounced),
    enabled: searchDebounced.trim().length > 0,
  });

  const { mutate: addCourse } = useMutation({
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

  const showPopover = isFocused && search.trim() !== "";

  return (
    <div
      onFocus={() => setIsFocused(true)}
      onBlurCapture={(e) => {
        if (!e.currentTarget.contains(e.relatedTarget)) {
          setIsFocused(false);
        }
      }}
      className="relative z-50"
    >
      <SearchInput
        value={search}
        onValueChange={setSearch}
        placeholder="Search and add offered courses (e.g., Computer Networks, CSE 3101)..."
        searching={isSearching}
        showClearButton={false}
        className="[&>div:last-child]:right-3 [&>input]:h-10 [&>input]:pl-9 [&>input]:text-sm [&>input]:shadow-sm [&>svg]:size-4 [&>svg:first-child]:left-3 [&>svg:last-child]:right-3 [&>svg:last-child]:size-5"
      />

      <Card
        className={cn(
          "absolute mt-2 max-h-80 w-full pt-0 pb-0",
          "scrollbar-thin scrollbar-thumb-accent overflow-y-auto",
          "rounded-lg bg-popover/95 shadow-xl backdrop-blur-md",
          "border border-muted/40",
          "animate-in duration-200",
          showPopover
            ? "translate-y-0 opacity-100"
            : "pointer-events-none -translate-y-4 opacity-0",
        )}
      >
        <div className="flex flex-col">
          {isSearching ? (
            <div className="flex h-20 items-center justify-center text-sm text-muted-foreground">
              <Loader2Icon className="mr-2 size-4 animate-spin" />
              Searching...
            </div>
          ) : searchResults && searchResults.length > 0 ? (
            searchResults.map((course) => (
              <SearchResultItem key={course.classID} course={course} onAdd={addCourse} />
            ))
          ) : (
            <div className="flex h-20 items-center justify-center text-sm text-muted-foreground">
              No matching offered courses found in the database.
            </div>
          )}
        </div>
      </Card>
    </div>
  );
}

function SearchResultItem({ course, onAdd }: { course: Course; onAdd: (classId: string) => void }) {
  return (
    <button
      type="button"
      key={course.classID}
      tabIndex={0}
      onClick={() => onAdd(course.classID)}
      onKeyDown={(e) => {
        if (e.key === "Enter") {
          e.preventDefault();
          onAdd(course.classID);
        }
      }}
      className={cn(
        "flex cursor-pointer items-center justify-between border-b px-4 py-3 transition-colors last:border-0 hover:bg-accent/60",
        "outline-none focus-visible:rounded focus-visible:ring-3 focus-visible:ring-ring/50 focus-visible:ring-inset",
      )}
    >
      <div className="flex flex-col gap-1">
        <div className="flex items-center gap-2 text-sm font-semibold">
          {course.courseCode && (
            <Badge variant="secondary" className="px-1.5 py-0.5 font-mono text-[9px]">
              {course.courseCode}
            </Badge>
          )}
          <span>{course.courseTitle}</span>
          <span className="text-xs font-normal text-muted-foreground">[{course.section}]</span>
        </div>
        <div className="flex items-center gap-4 text-[11px] text-muted-foreground">
          <span className="flex items-center gap-1">
            <UserIcon className="size-3" /> {course.faculty}
          </span>
          <span className="flex items-center gap-1">
            <CalendarIcon className="size-3" /> {course.schedules.map((s) => s.day).join(", ")}
          </span>
          <span className="flex items-center gap-1">
            <ClockIcon className="size-3" />
            {[...new Set(course.schedules.map((s) => `${s.startTime} - ${s.endTime}`))].join(", ")}
          </span>
        </div>
      </div>
    </button>
  );
}

function EmptyState() {
  return (
    <div className="flex min-h-[300px] flex-1 flex-col items-center justify-center py-20 text-center">
      <div className="mb-4 rounded-full bg-muted p-4">
        <CalendarIcon className="size-10 text-muted-foreground/60" />
      </div>
      <h3 className="text-lg font-bold">Your routine is empty</h3>
      <p className="mt-1 max-w-sm text-xs leading-relaxed text-muted-foreground">
        Search for your offered courses in the bar above or import an offered Excel sheet to
        construct your weekly timeline.
      </p>
    </div>
  );
}

interface DayScheduleTimelineProps {
  routineByDay: Record<string, Course[]>;
  onRemoveCourse: (classId: string) => void;
}

function DayScheduleTimeline({ routineByDay, onRemoveCourse }: DayScheduleTimelineProps) {
  return (
    <div className="flex flex-col gap-8 pb-12">
      {DAYS.map((day) => {
        const dayCourses = routineByDay[day];
        if (!dayCourses || dayCourses.length === 0) return null;

        return (
          <div
            key={day}
            className="flex animate-in flex-col gap-6 border-b pb-6 duration-300 fade-in last:border-0 lg:flex-row"
          >
            {/* Sticky Side Day Column */}
            <div className="flex shrink-0 items-baseline justify-between gap-1 lg:w-36 lg:flex-col lg:items-start">
              <h2 className="text-lg font-black tracking-tight text-foreground">{day}</h2>
              <Badge variant="secondary" className="px-2 py-0.5 text-[10px] font-bold">
                {dayCourses.length} {dayCourses.length === 1 ? "Class" : "Classes"}
              </Badge>
            </div>

            {/* Day Grid Courses */}
            <div className="grid flex-1 grid-cols-1 gap-4 md:grid-cols-2 xl:grid-cols-3">
              {dayCourses.flatMap((course) =>
                course.schedules
                  .filter((schedule) => schedule.day === day)
                  .map((schedule) => (
                    <CourseCard
                      key={`${course.classID}-${schedule.day}-${schedule.startTime}`}
                      course={course}
                      schedule={schedule}
                      onRemoveCourse={onRemoveCourse}
                    />
                  )),
              )}
            </div>
          </div>
        );
      })}
    </div>
  );
}

interface CourseCardProps {
  course: Course;
  schedule: Schedule;
  onRemoveCourse: (classId: string) => void;
}

function CourseCard({ course, schedule, onRemoveCourse }: CourseCardProps) {
  const status = getCourseStatus(schedule);
  const isLab = schedule.type.toLowerCase().includes("lab");

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
            <Badge variant="secondary" className="px-1 py-0 font-mono text-[9px] opacity-80">
              {course.courseCode}
            </Badge>
          )}
          <Badge variant="outline" className="-ml-1 px-2 py-0 text-[0.65rem] font-medium">
            {schedule.type}
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
          <span className="mt-0.5 block text-xs font-normal text-muted-foreground">
            Section {course.section}
          </span>
        </CardTitle>
      </CardHeader>

      <CardContent className="flex flex-col gap-2 pt-0">
        <div className="flex items-center gap-2 text-xs text-muted-foreground">
          <ClockIcon className="size-3.5 text-foreground/60" />
          <span className="font-medium text-foreground/90">
            {schedule.startTime} - {schedule.endTime}
          </span>
        </div>
        <div className="flex items-center gap-2 text-xs text-muted-foreground">
          <MapPinIcon className="size-3.5 text-foreground/60" />
          <span className="truncate">
            Room {schedule.room} ({course.department})
          </span>
        </div>
        <div className="flex items-center gap-2 text-xs text-muted-foreground">
          <UserIcon className="size-3.5 text-foreground/60" />
          <span className="truncate">{course.faculty}</span>
        </div>
      </CardContent>

      {/* Interactive Trash Trigger */}
      <Button
        variant="destructive"
        size="icon"
        className="absolute right-3 bottom-3 h-8 w-8 rounded-md opacity-0 transition-opacity duration-200 group-hover:opacity-100"
        onClick={() => onRemoveCourse(course.classID)}
      >
        <Trash2Icon className="size-4" />
      </Button>
    </Card>
  );
}
