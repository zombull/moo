// npm install --save-dev del gulp gulp-jshint gulp-concat gulp-inject gulp-rename gulp-replace gulp-uglify gulp-ng-annotate gulp-minify-css gulp-minify-html gulp-nodemon gulp-livereload

var del = require('del'),
    crypto = require('crypto'),
    fs = require('fs'),
    gulp = require('gulp'),
    jshint = require('gulp-jshint'),
    concat = require('gulp-concat'),
    inject = require('gulp-inject'),
    rename = require('gulp-rename'),
    replace = require('gulp-replace'),
    uglify = require('gulp-uglify'),
    ngAnnotate = require('gulp-ng-annotate'),
    minifyCss = require('gulp-minify-css'),
    minifyHTML = require('gulp-minify-html'),
    minifyInline = require('gulp-minify-inline'),
    nodemon = require('gulp-nodemon'),
    livereload = require('gulp-livereload');

var source = {
    js: ['moonboard/js/app.js', 'common/js/**/*.js', 'moonboard/js/**/*.js'],
    css: ['common/css/**/*.css', 'moonboard/css/**/*.css'],
    html: ['common/html/**/*.html', 'moonboard/html/**/*.html'],
    index: 'moonboard/index.html',
    substorage: 'moonboard/substorage.html',
    images: ['common/img/**/*', 'moonboard/img/**/*'],
    fonts: ['common/css/ocr/*.*', '!common/css/ocr/*.css',
            'common/css/universalia/*.*', '!common/css/universalia/*.css',
            'common/css/v5_bloques/*.*', '!common/css/v5_bloques/*.css'],
};

var destination = {
    moo: 'release/moo',
}

gulp.task('clean', function() {
    return del('release', {force: true});
});

// Run JS through jshint to find issues
gulp.task('jshint', function() {
    var jsHintOptions = {
        eqnull: true,
        "-W018": true,
        "-W041": false
    };
    return gulp.src(source.js)
        .pipe(jshint(jsHintOptions))
        .pipe(jshint.reporter('default'))
        .pipe(jshint.reporter('fail'));
});

// Write the checksums into database.js so that we don't have to make a REST request at runtime.
// TODO: update this to read from the checksums directory.
// gulp.task('checksums', function() {
//     var localdb = require('./utils/localdb');
//     return gulp.src('public/js/services/database.js')
//         .pipe(replace(/var\schecksums\s=.*;/, 'var checksums = ' + JSON.stringify(localdb.checksums).replace(/\"/g, "'") + ';'))
//         .pipe(gulp.dest('public/js/services'));
// });

function gulpYear(year) {
    var checksum = function(filepath, file) {
        filepath = filepath + '?version=' + crypto.createHash('md5').update(file.contents.toString('utf8')).digest('hex');
        return inject.transform.apply(inject.transform, arguments);
    }

    var sums = fs.readFileSync("moonboard/moonboard" + year + "md5.json", "utf8");

    var server = destination.moo + '/moonboard' + year;

    var html = gulp.src(source.html)
        .pipe(minifyHTML({empty: true}))                         // Minify HTML.  The empty option tells minifyHTML to keep empty attributes.
        .pipe(gulp.dest(server + '/html'));

    var css = gulp.src(source.css)
        .pipe(concat('moon.css'))                               // Concatenate everything into a single JS file.
        .pipe(gulp.dest(server + '/css'))                       // Save concatenated file before minification.
        .pipe(minifyCss())
        .pipe(rename({extname: ".min.css"}))                    // Rename the stream
        .pipe(gulp.dest(server + '/css'));

    var js = gulp.src(source.js)
        .pipe(concat('moon.js'))                                // Concatenate everything into a single JS file.
        .pipe(replace('SCHEMA_CHECKSUMS', sums))
        .pipe(replace('ZZZZ', year))                            // Replace year in subdomain references
        .pipe(replace('xyz:3000', 'xyz'))                       // Strip port off subdomain references
        .pipe(ngAnnotate({add: true, single_quotes: true}))     // Annotate angular code
        .pipe(gulp.dest(server + '/js'))                        // Save concatenated and annotated file before minification.
        .pipe(rename({extname: ".min.js"}))                     // Rename the stream
        .pipe(uglify())
        .pipe(gulp.dest(server + '/js'));

    gulp.src(source.fonts)
        .pipe(gulp.dest(server + '/css'));

    gulp.src(source.images)
        .pipe(gulp.dest(server + '/img'));

    gulp.src(source.substorage)
        .pipe(replace('ZZZZ', year))                            // Replace year in subdomain references
        .pipe(replace('xyz:3000', 'xyz'))                       // Strip port off subdomain references
        .pipe(minifyHTML({empty: true}))                        // Minify HTML.  The empty option tells minifyHTML to keep empty attributes.
        .pipe(minifyInline())
        .pipe(gulp.dest(server));

    gulp.src(source.index)
        .pipe(gulp.dest(server))                                // Necessary to set the path so injection works correctly.
        .pipe(replace('ZZZZ', year))                            // Replace year in subdomain references
        .pipe(replace('xyz:3000', 'xyz'))                       // Strip port off subdomain references
        .pipe(replace('ng-app', 'ng-strict-di ng-app'))
        .pipe(inject(css, {relative: true, addPrefix: 'static', transform: checksum}))
        .pipe(inject(js, {relative: true, addPrefix: 'static', transform: checksum}))
        .pipe(minifyHTML({empty: true}))                        // Minify HTML.  The empty option tells minifyHTML to keep empty attributes.
        .pipe(gulp.dest(server));
}

// This is a rather large task, but it makes sense because we're injecting CSS and JS into index.html.
// The overall amount of code is relatively small so it's not like it's taking a huge amount of time.
// The alternative would be to duplicate generation of the CSS and JS paths
gulp.task('server', /*['checksums'],*/ function() {
    gulpYear('2016');
    gulpYear('2017');
});

gulp.task('release', ['jshint', 'server']);

gulp.task('default', [
  'release'
]);
