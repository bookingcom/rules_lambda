set -e

TEST_PACKAGE="$(echo ${TEST_TARGET} | sed -e 's/:.*$//' -e 's@//@@')"
declare -r DATA_DIR="${TEST_SRCDIR}/${TEST_WORKSPACE}/${TEST_PACKAGE}"

for pkg in test_tar.tar test_zip.zip ; do
  ls -l "${DATA_DIR}/$pkg"
done

echo "PASS"
