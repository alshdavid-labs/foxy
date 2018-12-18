const request = require('request')
const fs = require('fs')
const path = require('path')

if (process.argv[2] === '--pre') {

    const available = {
        win32: {
            filename: 'foxy',
            url: 'https://github.com/alshdavid/foxy/raw/master/bin/win32/foxy.exe'
        },
        linux: {
            filename: 'foxy',
            url: 'https://github.com/alshdavid/foxy/raw/master/bin/linux/linux'
        }, 
        darwin: {
            filename: 'foxy',
            url: 'https://github.com/alshdavid/foxy/raw/master/bin/linux/linux'
        }
    }

    const option = available[process.platform]

    if (!option) {
        process.exit(1)
    }

    const download = (url, dest, cb = () => {}) => {
        const file = fs.createWriteStream(dest);
        const sendReq = request.get(url);

        // verify response code
        sendReq.on('response', (response) => {
            if (response.statusCode !== 200) {
                return cb('Response status was ' + response.statusCode);
            }

            sendReq.pipe(file);
        });

        // close() is async, call cb after close completes
        file.on('finish', () => file.close(cb));

        // check for request errors
        sendReq.on('error', (err) => {
            fs.unlink(dest);
            return cb(err.message);
        });

        file.on('error', (err) => { // Handle errors
            fs.unlink(dest); // Delete the file async. (But we don't check the result)
            return cb(err.message);
        });
    };

    if (!fs.existsSync(path.join(__dirname, ".bin"))) {
        fs.mkdirSync(path.join(__dirname, ".bin"))
    }
    download(option.url, path.join('.bin', option.filename))
}

if (process.argv[2] === '--post') {
    if (process.platform === 'win32') {
        fs.renameSync(path.join(__dirname, ".bin", "foxy"), path.join(__dirname, ".bin", "foxy.exe"))
    }
}