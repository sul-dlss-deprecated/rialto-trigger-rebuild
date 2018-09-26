default: package

package:
	GOOS=linux go build -o trigger-rebuild
	zip lambda.zip trigger-rebuild
