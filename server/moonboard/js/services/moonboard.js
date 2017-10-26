/**
 * Service for interacting with the actual board.
*/
moon.factory('moonboard', function ($q, $document, database) {
    'use strict';

    var board = {
        holds: [
            { img: 1, loc: 'H7', dir: 'SE' },
            { img: 2, loc: 'J14', dir: 'NW' },
            { img: 3, loc: 'K7', dir: 'N' },
            { img: 4, loc: 'D8', dir: 'N' },
            { img: 5, loc: 'A16', dir: 'NW' },
            { img: 6, loc: 'F6', dir: 'E' },
            { img: 7, loc: 'K6', dir: 'N' },
            { img: 8, loc: 'C9', dir: 'W' },
            { img: 9, loc: 'A10', dir: 'SW' },
            { img: 10, loc: 'I8', dir: 'N' },
            { img: 11, loc: 'K12', dir: 'N' },
            { img: 12, loc: 'A11', dir: 'SE' },
            { img: 13, loc: 'B7', dir: 'S' },
            { img: 14, loc: 'D16', dir: 'N' },
            { img: 15, loc: 'K10', dir: 'S' },
            { img: 16, loc: 'G16', dir: 'N' },
            { img: 17, loc: 'F8', dir: 'N' },
            { img: 18, loc: 'F15', dir: 'N' },
            { img: 19, loc: 'G7', dir: 'SW' },
            { img: 20, loc: 'H15', dir: 'NE' },
            { img: 21, loc: 'C12', dir: 'N' },
            { img: 22, loc: 'D6', dir: 'N' },
            { img: 23, loc: 'D3', dir: 'S' },
            { img: 24, loc: 'D13', dir: 'N' },
            { img: 25, loc: 'J11', dir: 'N' },
            { img: 26, loc: 'A13', dir: 'N' },
            { img: 27, loc: 'G11', dir: 'E' },
            { img: 28, loc: 'H18', dir: 'N' },
            { img: 29, loc: 'B4', dir: 'SW' },
            { img: 30, loc: 'K8', dir: 'N' },
            { img: 31, loc: 'C15', dir: 'NW' },
            { img: 32, loc: 'H9', dir: 'E' },
            { img: 33, loc: 'D10', dir: 'NW' },
            { img: 34, loc: 'H14', dir: 'W' },
            { img: 35, loc: 'I5', dir: 'N' },
            { img: 36, loc: 'I12', dir: 'SW' },
            { img: 37, loc: 'K13', dir: 'N' },
            { img: 38, loc: 'C7', dir: 'N' },
            { img: 39, loc: 'C18', dir: 'N' },
            { img: 40, loc: 'F9', dir: 'N' },
            { img: 50, loc: 'C14', dir: 'N' },
            { img: 51, loc: 'D17', dir: 'N' },
            { img: 52, loc: 'D9', dir: 'NE' },
            { img: 53, loc: 'F7', dir: 'NW' },
            { img: 54, loc: 'F12', dir: 'E' },
            { img: 55, loc: 'G12', dir: 'NE' },
            { img: 56, loc: 'B11', dir: 'NW' },
            { img: 57, loc: 'J10', dir: 'NE' },
            { img: 58, loc: 'J2', dir: 'SE' },
            { img: 59, loc: 'E13', dir: 'N' },
            { img: 60, loc: 'I6', dir: 'NE' },
            { img: 61, loc: 'J9', dir: 'SE' },
            { img: 62, loc: 'F14', dir: 'NW' },
            { img: 63, loc: 'I13', dir: 'E' },
            { img: 64, loc: 'E10', dir: 'NW' },
            { img: 65, loc: 'F10', dir: 'NE' },
            { img: 66, loc: 'E15', dir: 'NW' },
            { img: 67, loc: 'B8', dir: 'N' },
            { img: 68, loc: 'A12', dir: 'E' },
            { img: 69, loc: 'I16', dir: 'NE' },
            { img: 70, loc: 'I11', dir: 'N' },
            { img: 71, loc: 'B16', dir: 'NW' },
            { img: 72, loc: 'E11', dir: 'N' },
            { img: 73, loc: 'H11', dir: 'W' },
            { img: 74, loc: 'E7', dir: 'S' },
            { img: 75, loc: 'D12', dir: 'N' },
            { img: 76, loc: 'J8', dir: 'N' },
            { img: 77, loc: 'B13', dir: 'NW' },
            { img: 78, loc: 'B9', dir: 'NE' },
            { img: 79, loc: 'C10', dir: 'NE' },
            { img: 80, loc: 'B3', dir: 'SW' },
            { img: 81, loc: 'G2', dir: 'N' },
            { img: 82, loc: 'G18', dir: 'W' },
            { img: 83, loc: 'I4', dir: 'NE' },
            { img: 84, loc: 'K11', dir: 'NW' },
            { img: 85, loc: 'A5', dir: 'N' },
            { img: 86, loc: 'K5', dir: 'N' },
            { img: 87, loc: 'K18', dir: 'W' },
            { img: 88, loc: 'G8', dir: 'N' },
            { img: 89, loc: 'F5', dir: 'N' },
            { img: 90, loc: 'G13', dir: 'N' },
            { img: 91, loc: 'E18', dir: 'N' },
            { img: 92, loc: 'J6', dir: 'S' },
            { img: 93, loc: 'D14', dir: 'N' },
            { img: 94, loc: 'C11', dir: 'W' },
            { img: 95, loc: 'C6', dir: 'S' },
            { img: 96, loc: 'F16', dir: 'S' },
            { img: 97, loc: 'D5', dir: 'NW' },
            { img: 98, loc: 'A15', dir: 'N' },
            { img: 99, loc: 'B18', dir: 'SE' },
            { img: 100, loc: 'H16', dir: 'N' },
            { img: 101, loc: 'B15', dir: 'N' },
            { img: 102, loc: 'J12', dir: 'NE' },
            { img: 103, loc: 'J13', dir: 'N' },
            { img: 104, loc: 'K16', dir: 'N' },
            { img: 105, loc: 'F13', dir: 'NW' },
            { img: 106, loc: 'E16', dir: 'NW' },
            { img: 107, loc: 'I7', dir: 'NE' },
            { img: 108, loc: 'I15', dir: 'NW' },
            { img: 109, loc: 'I9', dir: 'SE' },
            { img: 110, loc: 'E12', dir: 'NE' },
            { img: 111, loc: 'H5', dir: 'NW' },
            { img: 112, loc: 'G15', dir: 'NW' },
            { img: 113, loc: 'J7', dir: 'N' },
            { img: 114, loc: 'H12', dir: 'NW' },
            { img: 115, loc: 'G17', dir: 'N' },
            { img: 116, loc: 'E9', dir: 'NE' },
            { img: 117, loc: 'J16', dir: 'E' },
            { img: 118, loc: 'F11', dir: 'NE' },
            { img: 119, loc: 'D11', dir: 'SW' },
            { img: 120, loc: 'I10', dir: 'N' },
            { img: 121, loc: 'K9', dir: 'N' },
            { img: 122, loc: 'E8', dir: 'N' },
            { img: 123, loc: 'A14', dir: 'NW' },
            { img: 124, loc: 'I14', dir: 'NW' },
            { img: 125, loc: 'C5', dir: 'N' },
            { img: 126, loc: 'D15', dir: 'NW' },
            { img: 127, loc: 'E14', dir: 'E' },
            { img: 128, loc: 'G9', dir: 'NE' },
            { img: 129, loc: 'E6', dir: 'NW' },
            { img: 130, loc: 'J5', dir: 'NW' },
            { img: 131, loc: 'H8', dir: 'NE' },
            { img: 132, loc: 'I18', dir: 'NE' },
            { img: 133, loc: 'A9', dir: 'NW' },
            { img: 134, loc: 'G6', dir: 'SW' },
            { img: 135, loc: 'C8', dir: 'NW' },
            { img: 136, loc: 'D18', dir: 'N' },
            { img: 137, loc: 'G14', dir: 'E' },
            { img: 138, loc: 'C13', dir: 'NW' },
            { img: 139, loc: 'A18', dir: 'N' },
            { img: 140, loc: 'H10', dir: 'NE' },
            { img: 141, loc: 'G4', dir: 'N' },
            { img: 142, loc: 'B12', dir: 'SE' },
            { img: 143, loc: 'C16', dir: 'N' },
            { img: 144, loc: 'K14', dir: 'NE' },
            { img: 145, loc: 'G10', dir: 'NE' },
            { img: 146, loc: 'D7', dir: 'S' },
            { img: 147, loc: 'B6', dir: 'NW' },
            { img: 148, loc: 'B10', dir: 'SE' },
            { img: 149, loc: 'H13', dir: 'SW' }
        ],

        xcoords: {
            A: 62,
            B: 94,
            C: 128,
            D: 160,
            E: 193,
            F: 226,
            G: 259,
            H: 292,
            I: 325,
            J: 357,
            K: 390,
        },
        ycoords: {
            '1': 616,
            '2': 582,
            '3': 549,
            '4': 516,
            '5': 484,
            '6': 451,
            '7': 417,
            '8': 385,
            '9': 352,
            '10': 319,
            '11': 287,
            '12': 254,
            '13': 221,
            '14': 189,
            '15': 155,
            '16': 122,
            '17': 89,
            '18': 56,
        },
        rotation: {
            'N': 0,
            'NE': 45,
            'E': 90,
            'SE': 135,
            'S': 180,
            'SW': 225,
            'W': 270,
            'NW': 315,
        }
    }

    var canvas = null;
    var stage = null;
    var borders = [];
    var containers = null;

    var images = null;

    function drawBoard(loaded) {
        canvas = document.getElementById("problemCanvas");
        stage = new createjs.Stage(canvas);
        borders = [];
        containers = new Map;
    
        stage.addChild(new createjs.Bitmap('data:image/png;base64,'+images[0]));
        
        _.each(board.holds, function(hold) {
            var img = new Image();
            img.src = 'data:image/png;base64,'+images[hold.img];
            img.onload = function (event) {
                var container = new createjs.Container();
                container.regX = (this.width / 2);
                container.regY = (this.height / 2);
                container.x = board.xcoords[hold.loc.substring(0, 1)];
                container.y = board.ycoords[hold.loc.substring(1)];
                container.rotation = board.rotation[hold.dir];
                container.z = 1;
    
                var bmp = new createjs.Bitmap(this);
                container.addChild(bmp);
                stage.addChild(container);
                
                containers.set(hold.loc, container);
                if (containers.size == board.holds.length) {
                    stage.update();
                    loaded.resolve();
                }
            }
        });
    }

    function setHolds(holds, color)  {
        _.each(holds, function(loc, id) {
            var container = containers.get(loc);
            var border = new createjs.Shape();
            border.graphics.setStrokeStyle(4).beginStroke(color).drawCircle(container.regX, container.regY, 20);
            container.addChild(border);
            container.z = 2;
            borders.push(border);
        });
    }

    return {
        load: function() {
            var loaded = $q.defer();
            $document.ready(function() {
                if (images === null) {
                    database.images(function(img) {
                        images = img;
                        drawBoard(loaded)
                    });
                } else {
                    drawBoard(loaded)
                }
            });
            return loaded.promise;
        },
        set: function(holds) {
            _.each(borders, function(border) {
                border.parent.removeChild(border);
            });
            borders = [];

            setHolds(holds.s, "rgba(0,0,0,1)");
            setHolds(holds.i, "rgba(255,57,0,1)");
            setHolds(holds.f, "rgba(0,0,0,1)");

            // Sort the children so that the board (z==0) is at the bottom
            // of the canvas, holds without borders are in the middle, and
            // holds with borders are on top.  This way the borders (circles)
            // are layered on top of any holds not in use by the problem.
            stage.sortChildren(function(a, b) {
                return a.z - b.z;
            });
            stage.update();
        }
    };
});
