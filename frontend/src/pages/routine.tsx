import { Badge } from "@/components/ui/badge";
import { Button } from "@/components/ui/button";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { Input } from "@/components/ui/input";
import { useDebounce } from "@/hooks/use-debounce";
import { Course } from "@bindings/routine";
import {
  AddToUserRoutine,
  GetUserRoutine,
  ImportOfferedCourses,
  RemoveFromUserRoutine,
  SearchOfferedCourses,
} from "@bindings/routine/service";
import { useQuery, useQueryClient } from "@tanstack/react-query";
import { Dialogs } from "@wailsio/runtime";
import {
  CalendarIcon,
  ClockIcon,
  FileSpreadsheetIcon,
  Loader2Icon,
  MapPinIcon,
  SearchIcon,
  Trash2Icon,
  UserIcon,
} from "lucide-react";
import { useState } from "react";
import { toast } from "sonner";

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

  const routineQuery = useQuery({
    queryKey: ["routine"],
    queryFn: () => GetUserRoutine(),
  });

  const routine = routineQuery.data ?? [];
  const loading = routineQuery.isLoading;

  const [search, setSearch] = useState("");
  const searchDebounced = useDebounce(search, 300);

  const courseQuery = useQuery({
    queryKey: ["courses", searchDebounced],
    queryFn: () => SearchOfferedCourses(searchDebounced),
    enabled: searchDebounced.trim().length > 0,
  });

  const searchResults = courseQuery.data ?? [];
  const isSearching = courseQuery.isLoading;

  const handleAddCourse = async (classId: string) => {
    try {
      await AddToUserRoutine(classId);
      setSearch("");
      queryClient.invalidateQueries({ queryKey: ["routine"] });
      toast.success("Course added to routine");
    } catch (err) {
      toast.error("Failed to add course", {
        description: err instanceof Error ? err.message : String(err),
      });
    }
  };

  const handleRemoveCourse = async (classId: string) => {
    try {
      await RemoveFromUserRoutine(classId);
      toast.success("Course removed from routine");
      queryClient.invalidateQueries({ queryKey: ["routine"] });
    } catch (err) {
      toast.error("Failed to remove course", {
        description: err instanceof Error ? err.message : String(err),
      });
    }
  };

  const handleImportCourses = async () => {
    try {
      const path = await Dialogs.OpenFile({
        Title: "Select a file",
        Filters: [
          { DisplayName: "Excel files", Pattern: "*.xlsx;*.xls" },
          { DisplayName: "All files", Pattern: "*.*" },
        ],
      });
      if (!path) return;
      await ImportOfferedCourses(path);
      toast.success("Courses imported successfully");
    } catch (err) {
      toast.error("Failed to import courses", {
        description: err instanceof Error ? err.message : String(err),
      });
    }
  };

  const routineByDay = DAYS.reduce(
    (acc, day) => {
      acc[day] = routine.filter((c) => c.day === day);
      return acc;
    },
    {} as Record<string, Course[]>,
  );

  return (
    <div className="animate-in fade-in-10 mx-auto flex h-full flex-col gap-6 p-6 duration-200 lg:p-10">
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-3xl font-bold">Class Routine</h1>
          <p className="text-muted-foreground mt-1 text-sm">
            Manage your weekly class schedule
          </p>
        </div>

        <Button variant="outline" onClick={handleImportCourses}>
          <FileSpreadsheetIcon className="size-4" />
          Import Offered Courses
        </Button>
      </div>

      <div className="relative z-50">
        <div className="relative">
          <SearchIcon className="text-muted-foreground absolute top-1/2 left-5 size-5 -translate-1/2" />
          <Input
            value={search}
            onChange={(e) => setSearch(e.target.value)}
            placeholder="Search Course"
            className="h-10 pl-10 text-base shadow-sm"
          />
          {isSearching && (
            <Loader2Icon className="text-muted-foreground absolute top-1/2 right-0 size-5 -translate-1/2 animate-spin" />
          )}
        </div>

        {search !== "" && searchResults.length > 0 && (
          <Card className="scrollbar-thumb-accent fade-in-10 animate-in absolute mt-2 max-h-90 w-full scrollbar-thin overflow-y-auto scroll-smooth rounded p-0 shadow-lg duration-200">
            <div className="flex flex-col">
              {searchResults.map((course) => (
                <div
                  key={course.classID}
                  className="hover:bg-accent/50 flex cursor-pointer items-center justify-between border-b px-4 py-3 transition-colors last:border-0"
                  onClick={() => handleAddCourse(course.classID)}
                >
                  <div className="flex flex-col gap-1">
                    <div className="flex items-center gap-2 font-medium">
                      {course.courseCode && (
                        <Badge
                          variant="secondary"
                          className="font-mono text-[10px]"
                        >
                          {course.courseCode}
                        </Badge>
                      )}
                      <span>{course.courseTitle}</span>
                      <span className="text-muted-foreground text-sm font-normal">
                        [{course.section}]
                      </span>
                    </div>
                    <div className="text-muted-foreground flex items-center gap-3 text-sm">
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
                </div>
              ))}
            </div>
          </Card>
        )}

        {search !== "" && !isSearching && searchResults.length === 0 && (
          <Card className="text-muted-foreground animate-in fade-in-10 absolute mt-2 w-full rounded p-4 text-center text-sm shadow-lg duration-200">
            No courses found
          </Card>
        )}
      </div>

      {/*Weekly Schedule*/}
      {loading ? (
        <div className="flex flex-1 items-center justify-center">
          <Loader2Icon className="text-muted-foreground size-8 animate-spin" />
        </div>
      ) : routine.length === 0 ? (
        <div className="flex flex-1 flex-col items-center justify-center text-center">
          <CalendarIcon className="text-muted-foreground mb-4 size-16" />
          <h3 className="text-xl font-semibold">Your routine is empty</h3>
          <p className="text-muted-foreground mt-2">
            Search for your courses above to start building your weekly
            schedule.
          </p>
        </div>
      ) : (
        <div className="flex scrollbar-none flex-col gap-8 overflow-y-auto pb-12">
          {DAYS.map((day) => {
            const dayCources = routineByDay[day];
            if (!dayCources || dayCources.length === 0) return null;

            return (
              <div key={day} className="flex gap-8">
                {/*Day Label*/}
                <div className="w-32 shrink-0">
                  <h2 className="text-primary bg-background sticky top-0 text-xl font-bold">
                    {day}
                  </h2>
                  <p className="text-muted-foreground text-sm font-medium">
                    {dayCources.length} Class{dayCources.length === 1 && "es"}
                  </p>
                </div>

                {/*Course Cards*/}
                <div className="my-1 grid flex-1 grid-cols-2 gap-4 xl:grid-cols-3">
                  {dayCources.map((course) => (
                    <Card
                      key={course.classID}
                      className="group hover:ring-ring/50 relative overflow-hidden shadow-xs transition-all duration-150 hover:ring-3"
                    >
                      <div className="bg-primary absolute top-0 left-0 h-full w-1" />
                      <CardHeader>
                        <div className="flex items-center gap-2">
                          {course.courseCode && (
                            <Badge variant="secondary" className="text-xs">
                              {course.courseCode}
                            </Badge>
                          )}
                          <Badge variant="outline">{course.type}</Badge>
                        </div>
                        <CardTitle>
                          {course.courseTitle}{" "}
                          <span className="text-muted-foreground">
                            [{course.section}]
                          </span>
                        </CardTitle>
                      </CardHeader>
                      <CardContent className="flex flex-col gap-2">
                        <div className="text-muted-foreground flex items-center gap-2 text-sm">
                          <ClockIcon className="size-4" />
                          <span className="text-foreground font-medium">
                            {course.startTime} - {course.endTime}
                          </span>
                        </div>
                        <div className="text-muted-foreground flex items-center gap-2 text-sm">
                          <MapPinIcon className="size-4" />
                          <span>
                            {course.room} ({course.department})
                          </span>
                        </div>
                        <div className="text-muted-foreground flex items-center gap-2 text-sm">
                          <UserIcon className="size-4" />
                          <span>{course.faculty}</span>
                        </div>
                      </CardContent>

                      {/*Delete Button*/}
                      <Button
                        variant="destructive"
                        size="icon"
                        className="absolute top-3 right-3 opacity-0 transition-opacity group-hover:opacity-100"
                        onClick={() => handleRemoveCourse(course.classID)}
                      >
                        <Trash2Icon className="size-4" />
                      </Button>
                    </Card>
                  ))}
                </div>
              </div>
            );
          })}
        </div>
      )}
    </div>
  );
}
