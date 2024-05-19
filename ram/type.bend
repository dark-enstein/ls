type File:
  fd { i }
  Name { name }
  Type { type }

type FileType:
  Socket { path }
  Normal { path }
  Directory { path }

def type(file):
  match file:
    case FileType/Socket:
      return "Socket type"
    case FileType/Normal:
      return "Normal filetype"
    case FileType/Directory:
      return "Directory filetype"


def main:
  return type(FileType/Socket{ path:"./emo" })