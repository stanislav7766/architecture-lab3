'use strict';

const crypto = require('crypto'),
  fs = require('fs');

const separateInputsFiles = inputFiles =>
  inputFiles.map(fileName => fileName.split('.')[0]);

const getHash = (filePath, algorithm = 'md5') =>
  new Promise((resolve, reject) => {
    const md5Sum = crypto.createHash(algorithm);
    try {
      const rs = fs.ReadStream(filePath);
      rs.on('data', data => {
        md5Sum.update(data);
      });
      rs.on('end', () => resolve(md5Sum.digest('hex')));
      rs.on('error', err => reject(err));
    } catch (err) {
      return reject(err);
    }
  });

const writeFile = (filePath, data) =>
  new Promise((resolve, reject) => {
    try {
      fs.writeFile(filePath, data, err =>
        err ? reject(err) : resolve(filePath)
      );
    } catch (err) {
      reject(err);
    }
  });
const listDir = dirName =>
  new Promise((resolve, reject) => {
    try {
      fs.readdir(dirName, (err, files) => (err ? reject(err) : resolve(files)));
    } catch (err) {
      reject(err);
    }
  });

(async function main() {
  if (process.argv.slice(2).length < 2)
    throw new Error('must be addition args');
  const format = {
    txt: '.txt',
    res: '.res'
  };
  const [dirInput, dirOutput] = process.argv.slice(2);
  const inputFiles = separateInputsFiles(await listDir(dirInput));
  const hashPromises = inputFiles.map(fileName =>
    getHash(dirInput + fileName + format.txt)
  );
  const hashFiles = await Promise.all(hashPromises);
  const writePromises = hashFiles.map((hash, i) =>
    writeFile(dirOutput + inputFiles[i] + format.res, hash)
  );
  console.log(`Total number of processed files: ${writePromises.length}.`);

  await Promise.all(writePromises);
})();
