# CMAKE generated file: DO NOT EDIT!
# Generated by "Unix Makefiles" Generator, CMake Version 3.21

# Delete rule output on recipe failure.
.DELETE_ON_ERROR:

#=============================================================================
# Special targets provided by cmake.

# Disable implicit rules so canonical targets will work.
.SUFFIXES:

# Disable VCS-based implicit rules.
% : %,v

# Disable VCS-based implicit rules.
% : RCS/%

# Disable VCS-based implicit rules.
% : RCS/%,v

# Disable VCS-based implicit rules.
% : SCCS/s.%

# Disable VCS-based implicit rules.
% : s.%

.SUFFIXES: .hpux_make_needs_suffix_list

# Command-line flag to silence nested $(MAKE).
$(VERBOSE)MAKESILENT = -s

#Suppress display of executed commands.
$(VERBOSE).SILENT:

# A target that is always out of date.
cmake_force:
.PHONY : cmake_force

#=============================================================================
# Set environment variables for the build.

# The shell in which to execute make rules.
SHELL = /bin/sh

# The CMake executable.
CMAKE_COMMAND = /usr/local/Cellar/cmake/3.21.3/bin/cmake

# The command to remove a file.
RM = /usr/local/Cellar/cmake/3.21.3/bin/cmake -E rm -f

# Escaping for special characters.
EQUALS = =

# The top-level source directory on which CMake was run.
CMAKE_SOURCE_DIR = /Users/tarun/Desktop/graphics/cse580-p1

# The top-level build directory on which CMake was run.
CMAKE_BINARY_DIR = /Users/tarun/Desktop/graphics/cse580-p1

# Include any dependencies generated for this target.
include CMU462/deps/glfw/tests/CMakeFiles/events.dir/depend.make
# Include any dependencies generated by the compiler for this target.
include CMU462/deps/glfw/tests/CMakeFiles/events.dir/compiler_depend.make

# Include the progress variables for this target.
include CMU462/deps/glfw/tests/CMakeFiles/events.dir/progress.make

# Include the compile flags for this target's objects.
include CMU462/deps/glfw/tests/CMakeFiles/events.dir/flags.make

CMU462/deps/glfw/tests/CMakeFiles/events.dir/events.c.o: CMU462/deps/glfw/tests/CMakeFiles/events.dir/flags.make
CMU462/deps/glfw/tests/CMakeFiles/events.dir/events.c.o: CMU462/deps/glfw/tests/events.c
CMU462/deps/glfw/tests/CMakeFiles/events.dir/events.c.o: CMU462/deps/glfw/tests/CMakeFiles/events.dir/compiler_depend.ts
	@$(CMAKE_COMMAND) -E cmake_echo_color --switch=$(COLOR) --green --progress-dir=/Users/tarun/Desktop/graphics/cse580-p1/CMakeFiles --progress-num=$(CMAKE_PROGRESS_1) "Building C object CMU462/deps/glfw/tests/CMakeFiles/events.dir/events.c.o"
	cd /Users/tarun/Desktop/graphics/cse580-p1/CMU462/deps/glfw/tests && /Library/Developer/CommandLineTools/usr/bin/cc $(C_DEFINES) $(C_INCLUDES) $(C_FLAGS) -MD -MT CMU462/deps/glfw/tests/CMakeFiles/events.dir/events.c.o -MF CMakeFiles/events.dir/events.c.o.d -o CMakeFiles/events.dir/events.c.o -c /Users/tarun/Desktop/graphics/cse580-p1/CMU462/deps/glfw/tests/events.c

CMU462/deps/glfw/tests/CMakeFiles/events.dir/events.c.i: cmake_force
	@$(CMAKE_COMMAND) -E cmake_echo_color --switch=$(COLOR) --green "Preprocessing C source to CMakeFiles/events.dir/events.c.i"
	cd /Users/tarun/Desktop/graphics/cse580-p1/CMU462/deps/glfw/tests && /Library/Developer/CommandLineTools/usr/bin/cc $(C_DEFINES) $(C_INCLUDES) $(C_FLAGS) -E /Users/tarun/Desktop/graphics/cse580-p1/CMU462/deps/glfw/tests/events.c > CMakeFiles/events.dir/events.c.i

CMU462/deps/glfw/tests/CMakeFiles/events.dir/events.c.s: cmake_force
	@$(CMAKE_COMMAND) -E cmake_echo_color --switch=$(COLOR) --green "Compiling C source to assembly CMakeFiles/events.dir/events.c.s"
	cd /Users/tarun/Desktop/graphics/cse580-p1/CMU462/deps/glfw/tests && /Library/Developer/CommandLineTools/usr/bin/cc $(C_DEFINES) $(C_INCLUDES) $(C_FLAGS) -S /Users/tarun/Desktop/graphics/cse580-p1/CMU462/deps/glfw/tests/events.c -o CMakeFiles/events.dir/events.c.s

CMU462/deps/glfw/tests/CMakeFiles/events.dir/__/deps/getopt.c.o: CMU462/deps/glfw/tests/CMakeFiles/events.dir/flags.make
CMU462/deps/glfw/tests/CMakeFiles/events.dir/__/deps/getopt.c.o: CMU462/deps/glfw/deps/getopt.c
CMU462/deps/glfw/tests/CMakeFiles/events.dir/__/deps/getopt.c.o: CMU462/deps/glfw/tests/CMakeFiles/events.dir/compiler_depend.ts
	@$(CMAKE_COMMAND) -E cmake_echo_color --switch=$(COLOR) --green --progress-dir=/Users/tarun/Desktop/graphics/cse580-p1/CMakeFiles --progress-num=$(CMAKE_PROGRESS_2) "Building C object CMU462/deps/glfw/tests/CMakeFiles/events.dir/__/deps/getopt.c.o"
	cd /Users/tarun/Desktop/graphics/cse580-p1/CMU462/deps/glfw/tests && /Library/Developer/CommandLineTools/usr/bin/cc $(C_DEFINES) $(C_INCLUDES) $(C_FLAGS) -MD -MT CMU462/deps/glfw/tests/CMakeFiles/events.dir/__/deps/getopt.c.o -MF CMakeFiles/events.dir/__/deps/getopt.c.o.d -o CMakeFiles/events.dir/__/deps/getopt.c.o -c /Users/tarun/Desktop/graphics/cse580-p1/CMU462/deps/glfw/deps/getopt.c

CMU462/deps/glfw/tests/CMakeFiles/events.dir/__/deps/getopt.c.i: cmake_force
	@$(CMAKE_COMMAND) -E cmake_echo_color --switch=$(COLOR) --green "Preprocessing C source to CMakeFiles/events.dir/__/deps/getopt.c.i"
	cd /Users/tarun/Desktop/graphics/cse580-p1/CMU462/deps/glfw/tests && /Library/Developer/CommandLineTools/usr/bin/cc $(C_DEFINES) $(C_INCLUDES) $(C_FLAGS) -E /Users/tarun/Desktop/graphics/cse580-p1/CMU462/deps/glfw/deps/getopt.c > CMakeFiles/events.dir/__/deps/getopt.c.i

CMU462/deps/glfw/tests/CMakeFiles/events.dir/__/deps/getopt.c.s: cmake_force
	@$(CMAKE_COMMAND) -E cmake_echo_color --switch=$(COLOR) --green "Compiling C source to assembly CMakeFiles/events.dir/__/deps/getopt.c.s"
	cd /Users/tarun/Desktop/graphics/cse580-p1/CMU462/deps/glfw/tests && /Library/Developer/CommandLineTools/usr/bin/cc $(C_DEFINES) $(C_INCLUDES) $(C_FLAGS) -S /Users/tarun/Desktop/graphics/cse580-p1/CMU462/deps/glfw/deps/getopt.c -o CMakeFiles/events.dir/__/deps/getopt.c.s

CMU462/deps/glfw/tests/CMakeFiles/events.dir/__/deps/glad.c.o: CMU462/deps/glfw/tests/CMakeFiles/events.dir/flags.make
CMU462/deps/glfw/tests/CMakeFiles/events.dir/__/deps/glad.c.o: CMU462/deps/glfw/deps/glad.c
CMU462/deps/glfw/tests/CMakeFiles/events.dir/__/deps/glad.c.o: CMU462/deps/glfw/tests/CMakeFiles/events.dir/compiler_depend.ts
	@$(CMAKE_COMMAND) -E cmake_echo_color --switch=$(COLOR) --green --progress-dir=/Users/tarun/Desktop/graphics/cse580-p1/CMakeFiles --progress-num=$(CMAKE_PROGRESS_3) "Building C object CMU462/deps/glfw/tests/CMakeFiles/events.dir/__/deps/glad.c.o"
	cd /Users/tarun/Desktop/graphics/cse580-p1/CMU462/deps/glfw/tests && /Library/Developer/CommandLineTools/usr/bin/cc $(C_DEFINES) $(C_INCLUDES) $(C_FLAGS) -MD -MT CMU462/deps/glfw/tests/CMakeFiles/events.dir/__/deps/glad.c.o -MF CMakeFiles/events.dir/__/deps/glad.c.o.d -o CMakeFiles/events.dir/__/deps/glad.c.o -c /Users/tarun/Desktop/graphics/cse580-p1/CMU462/deps/glfw/deps/glad.c

CMU462/deps/glfw/tests/CMakeFiles/events.dir/__/deps/glad.c.i: cmake_force
	@$(CMAKE_COMMAND) -E cmake_echo_color --switch=$(COLOR) --green "Preprocessing C source to CMakeFiles/events.dir/__/deps/glad.c.i"
	cd /Users/tarun/Desktop/graphics/cse580-p1/CMU462/deps/glfw/tests && /Library/Developer/CommandLineTools/usr/bin/cc $(C_DEFINES) $(C_INCLUDES) $(C_FLAGS) -E /Users/tarun/Desktop/graphics/cse580-p1/CMU462/deps/glfw/deps/glad.c > CMakeFiles/events.dir/__/deps/glad.c.i

CMU462/deps/glfw/tests/CMakeFiles/events.dir/__/deps/glad.c.s: cmake_force
	@$(CMAKE_COMMAND) -E cmake_echo_color --switch=$(COLOR) --green "Compiling C source to assembly CMakeFiles/events.dir/__/deps/glad.c.s"
	cd /Users/tarun/Desktop/graphics/cse580-p1/CMU462/deps/glfw/tests && /Library/Developer/CommandLineTools/usr/bin/cc $(C_DEFINES) $(C_INCLUDES) $(C_FLAGS) -S /Users/tarun/Desktop/graphics/cse580-p1/CMU462/deps/glfw/deps/glad.c -o CMakeFiles/events.dir/__/deps/glad.c.s

# Object files for target events
events_OBJECTS = \
"CMakeFiles/events.dir/events.c.o" \
"CMakeFiles/events.dir/__/deps/getopt.c.o" \
"CMakeFiles/events.dir/__/deps/glad.c.o"

# External object files for target events
events_EXTERNAL_OBJECTS =

CMU462/deps/glfw/tests/events: CMU462/deps/glfw/tests/CMakeFiles/events.dir/events.c.o
CMU462/deps/glfw/tests/events: CMU462/deps/glfw/tests/CMakeFiles/events.dir/__/deps/getopt.c.o
CMU462/deps/glfw/tests/events: CMU462/deps/glfw/tests/CMakeFiles/events.dir/__/deps/glad.c.o
CMU462/deps/glfw/tests/events: CMU462/deps/glfw/tests/CMakeFiles/events.dir/build.make
CMU462/deps/glfw/tests/events: CMU462/deps/glfw/src/libglfw3.a
CMU462/deps/glfw/tests/events: CMU462/deps/glfw/tests/CMakeFiles/events.dir/link.txt
	@$(CMAKE_COMMAND) -E cmake_echo_color --switch=$(COLOR) --green --bold --progress-dir=/Users/tarun/Desktop/graphics/cse580-p1/CMakeFiles --progress-num=$(CMAKE_PROGRESS_4) "Linking C executable events"
	cd /Users/tarun/Desktop/graphics/cse580-p1/CMU462/deps/glfw/tests && $(CMAKE_COMMAND) -E cmake_link_script CMakeFiles/events.dir/link.txt --verbose=$(VERBOSE)

# Rule to build all files generated by this target.
CMU462/deps/glfw/tests/CMakeFiles/events.dir/build: CMU462/deps/glfw/tests/events
.PHONY : CMU462/deps/glfw/tests/CMakeFiles/events.dir/build

CMU462/deps/glfw/tests/CMakeFiles/events.dir/clean:
	cd /Users/tarun/Desktop/graphics/cse580-p1/CMU462/deps/glfw/tests && $(CMAKE_COMMAND) -P CMakeFiles/events.dir/cmake_clean.cmake
.PHONY : CMU462/deps/glfw/tests/CMakeFiles/events.dir/clean

CMU462/deps/glfw/tests/CMakeFiles/events.dir/depend:
	cd /Users/tarun/Desktop/graphics/cse580-p1 && $(CMAKE_COMMAND) -E cmake_depends "Unix Makefiles" /Users/tarun/Desktop/graphics/cse580-p1 /Users/tarun/Desktop/graphics/cse580-p1/CMU462/deps/glfw/tests /Users/tarun/Desktop/graphics/cse580-p1 /Users/tarun/Desktop/graphics/cse580-p1/CMU462/deps/glfw/tests /Users/tarun/Desktop/graphics/cse580-p1/CMU462/deps/glfw/tests/CMakeFiles/events.dir/DependInfo.cmake --color=$(COLOR)
.PHONY : CMU462/deps/glfw/tests/CMakeFiles/events.dir/depend

