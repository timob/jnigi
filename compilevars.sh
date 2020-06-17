jdk_path=$1

if [[ "$jdk_path" == "" ]]; then
	echo "jdk path not given"
	echo "usage: source ./compilevars.sh <jdk path>"
else

  if [[ "$OSTYPE" == "linux-gnu"* ]]; then
	  osdir=linux
  elif [[ "$OSTYPE" == "darwin"* ]]; then
  	  osdir=darwin
  elif [[ "$OSTYPE" == "msys" ]]; then
  	  osdir=win32
  fi

  export CGO_CFLAGS="-I${jdk_path}/include -I${jdk_path}/include/${osdir}"

  unset jdk_path
  unset osdir

fi
