namespace go compilationInterface

struct Command
{
	1:string commandLine,
	2:string id,
	3:string workingDir
}

typedef i32 int
service CompilationService
{
		int executeCommand(1:Command command),
}