namespace go compilationInterface

struct Command
{
	1:string commandLine,
	3:string id
}

typedef i32 int
service CompilationService
{
		int executeCommand(1:Command command),
}