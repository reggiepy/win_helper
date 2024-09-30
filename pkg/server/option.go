package server

type Option func(s *Server) error

func WithBasePath(basePath string) Option {
	return func(s *Server) error {
		s.BasePath = basePath
		return nil
	}
}

func WithSId(id string) Option {
	return func(s *Server) error {
		s.SId = id
		return nil
	}
}

func WithSExecutable(executable string) Option {
	return func(s *Server) error {
		s.SExecutable = executable
		return nil
	}
}

func WithSName(name string) Option {
	return func(s *Server) error {
		s.SName = name
		return nil
	}
}

func WithSDescription(description string) Option {
	return func(s *Server) error {
		s.SDescription = description
		return nil
	}
}

func WithSStartMode(startMode string) Option {
	return func(s *Server) error {
		s.SStartMode = startMode
		return nil
	}
}

func WithSDepends(depends []string) Option {
	return func(s *Server) error {
		s.SDepends = depends
		return nil
	}
}

func WithSLogPath(logPath string) Option {
	return func(s *Server) error {
		s.SLogPath = logPath
		return nil
	}
}

func WithSArguments(arguments string) Option {
	return func(s *Server) error {
		s.SArguments = arguments
		return nil
	}
}

func WithSStartArguments(startArguments string) Option {
	return func(s *Server) error {
		s.SStartArguments = startArguments
		return nil
	}
}

func WithSStopExecutable(stopExecutable string) Option {
	return func(s *Server) error {
		s.SStopExecutable = stopExecutable
		return nil
	}
}

func WithSStopArguments(stopArguments string) Option {
	return func(s *Server) error {
		s.SStopArguments = stopArguments
		return nil
	}
}

func WithSEnv(env []string) Option {
	return func(s *Server) error {
		s.SEnv = env
		return nil
	}
}

func WithSFailure(failure string) Option {
	return func(s *Server) error {
		s.SFailure = failure
		return nil
	}
}

func WithSWorkingDirectory(workingDirectory string) Option {
	return func(s *Server) error {
		s.SWorkingDirectory = workingDirectory
		return nil
	}
}

func WithSLogMode(logMode string) Option {
	return func(s *Server) error {
		s.SLogMode = logMode
		return nil
	}
}

func WithSLogPattern(logPattern string) Option {
	return func(s *Server) error {
		s.SLogPattern = logPattern
		return nil
	}
}

func WithSLogAutoRollAtTime(autoRollAtTime string) Option {
	return func(s *Server) error {
		s.SLogAutoRollAtTime = autoRollAtTime
		return nil
	}
}

func WithSLogSizeThreshold(sizeThreshold int) Option {
	return func(s *Server) error {
		s.SLogSizeThreshold = sizeThreshold
		return nil
	}
}

func WithSLogKeepFiles(logKeepFiles int) Option {
	return func(s *Server) error {
		s.SLogKeepFiles = logKeepFiles
		return nil
	}
}

func WithSForce(force bool) Option {
	return func(s *Server) error {
		s.sForce = force
		return nil
	}
}
